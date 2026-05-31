package utils

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"sentinel/config"
	"sentinel/internal/models"
	"sentinel/pkg/constants"
	"sentinel/pkg/logger"
	"sentinel/pkg/mail"

	"github.com/xuri/excelize/v2"
)

// This function converts a slice of strings to a single string with comma separated values
func ArrayToString(array []string) string {
	var str string
	for _, v := range array {
		str += v + ","
	}
	return str
}

// FindKeyValueInJson finds the value of a given key in a JSON string.
func FindKeyValueInJson(json string, key string) string {
	parameterList := strings.Split(json, ",")
	for _, v := range parameterList {
		// if value contains "username" string, split it by equal sign and get the second value
		if strings.Contains(v, "username") {
			username := strings.Split(v, ":")
			return username[1]
		}
	}
	return ""
}

// Filter to Logs by ignored error logs like ("AuthLoginFailed 1002: Geçersiz e-posta / kullanıcı adı veya şifre.")
func FilterChanges(changes []models.Log, ignored []string) []models.Log {
	var filteredChanges []models.Log

	for _, change := range changes {
		if len(ignored) > 0 {
			for _, v := range ignored {
				if change.Message != v {
					filteredChanges = append(filteredChanges, change)
				}
			}
		} else {
			filteredChanges = append(filteredChanges, change)
		}
	}
	return filteredChanges
}

func UrlToOptions(url string) (string, string, string, string, string, string) {

	options := strings.Split(url, "://")

	// split protocol and info
	protocol := options[0]
	info := options[1]

	// split info to username, password, host, port, db
	infoOptions := strings.Split(info, "@")

	// split username and password
	usernamePassword := infoOptions[0]
	hostPortDb := infoOptions[1]

	usernamePasswordOptions := strings.Split(usernamePassword, ":")
	username := usernamePasswordOptions[0]
	password := usernamePasswordOptions[1]

	// split host, port and db
	hostPortDbOptions := strings.Split(hostPortDb, "/")
	hostPort := hostPortDbOptions[0]
	db := hostPortDbOptions[1]

	// split host and port
	hostPortOptions := strings.Split(hostPort, ":")
	host := hostPortOptions[0]
	port := hostPortOptions[1]

	return protocol, username, password, host, port, db
}

// NormalizeDomain normalizes domain input to host:port format
// Supports formats:
//   - domain.com (adds default port 443)
//   - domain.com:443 (uses as-is)
//   - https://domain.com (parses URL, uses port 443)
//   - http://domain.com:8080 (parses URL, uses specified port)
//   - https://domain.com:8443 (parses URL, uses specified port)
func NormalizeDomain(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}

	// Trim spaces
	domain = strings.TrimSpace(domain)

	// Check if domain already has protocol
	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		// Parse URL
		parsedURL, err := url.Parse(domain)
		if err != nil {
			return "", fmt.Errorf("invalid URL format: %v", err)
		}

		host := parsedURL.Hostname()
		port := parsedURL.Port()

		// If port is not specified, use default based on scheme
		if port == "" {
			if parsedURL.Scheme == "https" {
				port = "443"
			} else {
				port = "80"
			}
		}

		return fmt.Sprintf("%s:%s", host, port), nil
	}

	// Check if domain already has port (format: domain:port)
	if strings.Contains(domain, ":") {
		// Validate that it has exactly one colon and port is numeric
		parts := strings.Split(domain, ":")
		if len(parts) == 2 {
			if _, err := strconv.Atoi(parts[1]); err == nil {
				// Valid domain:port format
				return domain, nil
			}
		}
	}

	// No protocol, no port - add default HTTPS port
	return fmt.Sprintf("%s:443", domain), nil
}

// Check Domain Certificate with context support
func CheckDomainCertificate(domain string, day int) (bool, *models.Log) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	return CheckDomainCertificateWithContext(ctx, domain, day)
}

// CheckDomainCertificateWithContext checks domain certificate with context support for timeout and cancellation
func CheckDomainCertificateWithContext(ctx context.Context, domain string, day int) (bool, *models.Log) {
	status := constants.CertStatusNotExpired

	if day <= 0 {
		day = constants.DefaultExpireDays
	}

	// Normalize domain to host:port format
	normalizedDomain, err := NormalizeDomain(domain)
	if err != nil {
		logger.CLogger.Error("Failed to normalize domain:", err)
		return false, nil
	}

	// false: certificate will not expire in 30 days
	// true: certificate will expire in 30 days
	logger.CLogger.Info("INFO: Checking certificate for " + normalizedDomain)

	// TCP connection to domain with context
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", normalizedDomain)
	if err != nil {
		if netErr, ok := err.(*net.OpError); ok && netErr.Op == "dial" {
			// DNS resolution error
			logger.CLogger.Error("Failed to establish TCP connection - DNS resolution error:", err)
		} else {
			// Other error
			logger.CLogger.Error("Failed to establish TCP connection:", err)
		}
		return false, nil
	}
	defer conn.Close()

	// Set deadline for connection based on context
	if deadline, ok := ctx.Deadline(); ok {
		conn.SetDeadline(deadline)
	}

	// TLS Handshake
	// x509: certificate signed by unknown authority
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         strings.Split(normalizedDomain, ":")[0],
		InsecureSkipVerify: true,
	})

	if err := tlsConn.Handshake(); err != nil {
		logger.CLogger.Error("TLS Handshake failed:", err)
		return false, nil
	}
	defer tlsConn.Close()

	// Get certificate info directly from TLS connection state
	// No need to read HTTP response body - just handshake is enough
	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		logger.CLogger.Error("No peer certificates found")
		return false, nil
	}

	cert := state.PeerCertificates[0]
	tempPort, _ := strconv.Atoi(strings.Split(normalizedDomain, ":")[1])
	tempOrganization := cert.Subject.Organization                      // Optimized line
	daysUntilExpiration := int(time.Until(cert.NotAfter).Hours() / 24) // Optimized line

	// Check if certificate is actually expired (current time > expiration time)
	actuallyExpired := time.Now().After(cert.NotAfter)

	// Check if certificate will expire within the specified days
	willExpireSoon := daysUntilExpiration < day

	// Set status based on expiration state
	if actuallyExpired {
		status = constants.CertStatusExpired
	} else if willExpireSoon {
		status = constants.CertStatusExpired // Warning: will expire soon
	}

	// if certifcate time gonna expire in 30 days add to logs
	return willExpireSoon, &models.Log{
		Version:            cert.Version,
		SerialNumber:       cert.SerialNumber.String(),
		Subject:            cert.Subject.String(),
		IssuerSubject:      cert.Issuer.String(),
		Domain:             strings.Split(normalizedDomain, ":")[0],
		Port:               tempPort,
		CommonName:         cert.Subject.CommonName,
		Organization:       ArrayToString(tempOrganization),
		IssuedOn:           cert.NotBefore,
		ExpiresOn:          cert.NotAfter,
		CertificateData:    string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		SubjectKeyID:       hex.EncodeToString(cert.SubjectKeyId),
		AuthorityKeyID:     hex.EncodeToString(cert.AuthorityKeyId),
		IsCA:               cert.IsCA,
		Issuer:             cert.Issuer.CommonName,
		IsExpired:          actuallyExpired,
		Message:            fmt.Sprintf("Certificate will expire in %d days.", daysUntilExpiration),
		Status:             status,
	}
}

// Decode Certificate Data (PEM format)
func DecodeCertificateData(certData string) (*x509.Certificate, error) {
	// Decode the PEM encoded certificate
	block, _ := pem.Decode([]byte(certData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Parse the certificate
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %v", err)
	}

	return cert, nil
}

// Excel File Creation Function
func SetChangesToExcel(changes []models.Log) *excelize.File {
	// Create a new spreadsheet
	f := excelize.NewFile()

	// Expire Style
	styleExpire, errExpire := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{constants.ExcelColorRed},
			Pattern: 1,
		},
	})
	if errExpire != nil {
		logger.CLogger.Error("ERROR: ", errExpire)
		f.Close()
		return nil
	}

	// Not Expire Style
	styleNotExpire, errNotExpire := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{constants.ExcelColorGreen},
			Pattern: 1,
		},
	})
	if errNotExpire != nil {
		logger.CLogger.Error("ERROR: ", errNotExpire)
		f.Close()
		return nil
	}

	// Time Out Style
	styleTimeOut, errTimeOut := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{constants.ExcelColorYellow},
			Pattern: 1,
		},
	})
	if errTimeOut != nil {
		logger.CLogger.Error("ERROR: ", errTimeOut)
		f.Close()
		return nil
	}

	// Change the name of the worksheet.
	f.SetSheetName("Sheet1", "Logs")

	// Set headers
	headers := []string{
		"Version", "Serial Number", "Subject", "Issuer Subject", "Domain",
		"Port", "Common Name", "Organization", "Issued On", "Expires On",
		"Certificate Data", "Signature Algorithm", "Subject Key ID",
		"Authority Key ID", "Is CA", "Issuer", "Is Expired", "Message",
	}

	for i, header := range headers {
		col := string(rune('A' + i))
		f.SetCellValue("Logs", col+"1", header)
	}

	// Set value of cells - optimized with batch operations
	for index, change := range changes {
		row := strconv.Itoa(index + 2)

		// Set all values for the row
		f.SetCellValue("Logs", "A"+row, change.Version)
		f.SetCellValue("Logs", "B"+row, change.SerialNumber)
		f.SetCellValue("Logs", "C"+row, change.Subject)
		f.SetCellValue("Logs", "D"+row, change.IssuerSubject)
		f.SetCellValue("Logs", "E"+row, change.Domain)
		f.SetCellValue("Logs", "F"+row, change.Port)
		f.SetCellValue("Logs", "G"+row, change.CommonName)
		f.SetCellValue("Logs", "H"+row, change.Organization)
		f.SetCellValue("Logs", "I"+row, change.IssuedOn)
		f.SetCellValue("Logs", "J"+row, change.ExpiresOn)
		f.SetCellValue("Logs", "K"+row, change.CertificateData)
		f.SetCellValue("Logs", "L"+row, change.SignatureAlgorithm)
		f.SetCellValue("Logs", "M"+row, change.SubjectKeyID)
		f.SetCellValue("Logs", "N"+row, change.AuthorityKeyID)
		f.SetCellValue("Logs", "O"+row, change.IsCA)
		f.SetCellValue("Logs", "P"+row, change.Issuer)
		f.SetCellValue("Logs", "Q"+row, change.IsExpired)
		f.SetCellValue("Logs", "R"+row, change.Message)

		// Apply style based on status
		var style int
		switch change.Status {
		case constants.CertStatusExpired:
			style = styleExpire
		case constants.CertStatusNotExpired:
			style = styleNotExpire
		default:
			style = styleTimeOut
		}
		f.SetCellStyle("Logs", "A"+row, "R"+row, style)
	}

	// Set active sheet
	f.SetActiveSheet(0)

	// Save spreadsheet
	if err := f.SaveAs("Logs.xlsx"); err != nil {
		logger.CLogger.Error("ERROR: ", err)
		f.Close()
		return nil
	}

	// return file for attachment
	return f
}

func SendMailWithAttachment(logs []models.Log, f *excelize.File) {
	// Check if mail is properly configured
	if config.C.Mail.Host == "" || config.C.Mail.Host == "HOST" {
		logger.CLogger.Warn("Mail is not configured - skipping email notification")
		return
	}

	// Get recipients from DataService (DB or static data based on config)
	toUsers := ToUsers
	ccUsers := CCUsers
	bccUsers := []string{}

	mailContent := &models.Mail{
		Sender:  config.C.Mail.FromMail,
		To:      toUsers,
		Cc:      ccUsers,
		Bcc:     bccUsers,
		Subject: config.C.App.Name + " Error Logs",
	}

	mail.SendMail(mailContent, logs, f)
}
