package helpers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"sentinel/logger"
	"sentinel/models"

	"github.com/xuri/excelize/v2"
)

func TimeFormatter(t time.Time) string {
	// That formats the given time to RFC3339 format (2023-01-16T13:50:56.910Z)
	return t.UTC().Format(time.RFC3339)
}

func ArrayToString(array []string) string {
	// That converts the given array to string
	var str string
	for _, v := range array {
		str += v + ","
	}
	return str
}

// Find key value in string json
func FindKeyValueInJson(json string, key string) string {
	// That finds the given key value in string json
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

// Check Domain Certificate
func CheckDomainCertificate(domain string, day int) (bool, *models.Log) {
	status := 0

	// false: certificate will not expire in 30 days
	// true: certificate will expire in 30 days
	logger.CLogger.Info("INFO: Checking certificate for " + domain)

	// TCP connection to domain
	conn, err := net.DialTimeout("tcp", domain, 10*time.Second)
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

	// TLS Handshake
	// x509: certificate signed by unknown authority
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         strings.Split(domain, ":")[0],
		InsecureSkipVerify: true,
	})

	if err := tlsConn.Handshake(); err != nil {
		logger.CLogger.Error("TLS Handshake failed:", err)
		return false, nil
	}

	// HTTP Request
	req := "GET / HTTP/1.1\r\nHost: " + strings.Split(domain, ":")[0] + "\r\n\r\n"
	if _, err := tlsConn.Write([]byte(req)); err != nil {
		logger.CLogger.Error("Failed to write HTTP request:", err)
		return false, nil
	}

	// HTTP Response
	var responseBuffer bytes.Buffer
	buf := make([]byte, 1024)
	for {
		n, err := tlsConn.Read(buf)
		if err != nil {
			if err != io.EOF {
				logger.CLogger.Error("Failed to read HTTP response:", err)
			}
			break
		}
		responseBuffer.Write(buf[:n])
	}

	// Certification Info is here
	cert := tlsConn.ConnectionState().PeerCertificates[0]
	tempPort, _ := strconv.Atoi(strings.Split(domain, ":")[1])
	tempOrganization := cert.Subject.Organization
	daysUntilExpiration := int(time.Until(cert.NotAfter).Hours() / 24) // Optimized line
	isExpired := daysUntilExpiration < day
	if isExpired {
		status = 1
	}

	// if certifcate time gonna expire in 30 days add to logs
	return isExpired, &models.Log{
		Version:            cert.Version,
		SerialNumber:       cert.SerialNumber.String(),
		Subject:            cert.Subject.String(),
		IssuerSubject:      cert.Issuer.String(),
		Domain:             strings.Split(domain, ":")[0],
		Port:               tempPort,
		CommonName:         cert.Subject.CommonName,
		Organization:       ArrayToString(tempOrganization),
		IssuedOn:           cert.NotBefore,
		ExpiresOn:          cert.NotAfter,
		CertificateData:    string(cert.Raw),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
		SubjectKeyID:       string(cert.SubjectKeyId),
		AuthorityKeyID:     string(cert.AuthorityKeyId),
		IsCA:               cert.IsCA,
		Issuer:             cert.Issuer.CommonName,
		IsExpired:          cert.NotAfter.Before(cert.NotBefore),
		Message:            fmt.Sprintf("Certificate will expire in %d days.", daysUntilExpiration),
		Status:             status,
	}
}

// Excel File Creation Function
func SetChangesToExcel(changes []models.Log) *excelize.File {

	// Create a new spreadsheet
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			logger.CLogger.Error("ERROR: ", err)
		}
	}()

	// Expire Style
	styleExpire, errExpire := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FF0000"},
			Pattern: 1,
		},
	})
	if errExpire != nil {
		logger.CLogger.Error("ERROR: ", errExpire)
		return nil
	}

	// Not Expire Style
	styleNotExpire, errNotExpire := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
	})
	if errNotExpire != nil {
		logger.CLogger.Error("ERROR: ", errNotExpire)
		return nil
	}

	// Time Out Style
	styleTimeOut, errTimeOut := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#FFFF00"},
			Pattern: 1,
		},
	})
	if errTimeOut != nil {
		logger.CLogger.Error("ERROR: ", errTimeOut)
		return nil
	}

	// Change the name of the worksheet.
	f.SetSheetName("Sheet1", "Logs")

	f.SetCellValue("Logs", "A1", "Version")
	f.SetCellValue("Logs", "B1", "Serial Number")
	f.SetCellValue("Logs", "C1", "Subject")
	f.SetCellValue("Logs", "D1", "Issuer Subject")
	f.SetCellValue("Logs", "E1", "Domain")
	f.SetCellValue("Logs", "F1", "Port")
	f.SetCellValue("Logs", "G1", "Common Name")
	f.SetCellValue("Logs", "H1", "Organization")
	f.SetCellValue("Logs", "I1", "Issued On")
	f.SetCellValue("Logs", "J1", "Expires On")
	f.SetCellValue("Logs", "K1", "Certificate Data")
	f.SetCellValue("Logs", "L1", "Signature Algorithm")
	f.SetCellValue("Logs", "M1", "Subject Key ID")
	f.SetCellValue("Logs", "N1", "Authority Key ID")
	f.SetCellValue("Logs", "O1", "Is CA")
	f.SetCellValue("Logs", "P1", "Issuer")
	f.SetCellValue("Logs", "Q1", "Is Expired")
	f.SetCellValue("Logs", "R1", "Message")

	// Set value of a cell.
	index := 2
	for _, change := range changes {
		f.SetCellValue("Logs", "A"+strconv.Itoa(index), change.Version)
		f.SetCellValue("Logs", "B"+strconv.Itoa(index), change.SerialNumber)
		f.SetCellValue("Logs", "C"+strconv.Itoa(index), change.Subject)
		f.SetCellValue("Logs", "D"+strconv.Itoa(index), change.IssuerSubject)
		f.SetCellValue("Logs", "E"+strconv.Itoa(index), change.Domain)
		f.SetCellValue("Logs", "F"+strconv.Itoa(index), change.Port)
		f.SetCellValue("Logs", "G"+strconv.Itoa(index), change.CommonName)
		f.SetCellValue("Logs", "H"+strconv.Itoa(index), change.Organization)
		f.SetCellValue("Logs", "I"+strconv.Itoa(index), change.IssuedOn)
		f.SetCellValue("Logs", "J"+strconv.Itoa(index), change.ExpiresOn)
		f.SetCellValue("Logs", "K"+strconv.Itoa(index), change.CertificateData)
		f.SetCellValue("Logs", "L"+strconv.Itoa(index), change.SignatureAlgorithm)
		f.SetCellValue("Logs", "M"+strconv.Itoa(index), change.SubjectKeyID)
		f.SetCellValue("Logs", "N"+strconv.Itoa(index), change.AuthorityKeyID)
		f.SetCellValue("Logs", "O"+strconv.Itoa(index), change.IsCA)
		f.SetCellValue("Logs", "P"+strconv.Itoa(index), change.Issuer)
		f.SetCellValue("Logs", "Q"+strconv.Itoa(index), change.IsExpired)
		f.SetCellValue("Logs", "R"+strconv.Itoa(index), change.Message)
		if change.Status == 1 {
			f.SetCellStyle("Logs", "A"+strconv.Itoa(index), "R"+strconv.Itoa(index), styleExpire)
		} else if change.Status == 0 {
			f.SetCellStyle("Logs", "A"+strconv.Itoa(index), "R"+strconv.Itoa(index), styleNotExpire)
		} else {
			f.SetCellStyle("Logs", "A"+strconv.Itoa(index), "R"+strconv.Itoa(index), styleTimeOut)
		}
		index++
	}

	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Save spreadsheet to the db
	if err := f.SaveAs("Logs.xlsx"); err != nil {
		logger.CLogger.Error("ERROR: ", err)
		return nil
	}

	// return file for attachment
	return f
}
