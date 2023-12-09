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
	daysUntilExpiration := int(time.Until(cert.NotAfter).Hours() / 24) // Optimized line
	isExpired := daysUntilExpiration < day

	// if certifcate time gonna expire in 30 days add to logs
	return isExpired, &models.Log{
		Version:            cert.Version,
		SerialNumber:       cert.SerialNumber.String(),
		Subject:            cert.Subject.String(),
		IssuerSubject:      cert.Issuer.String(),
		Domain:             strings.Split(domain, ":")[0],
		Port:               tempPort,
		CommonName:         cert.Subject.CommonName,
		Organization:       cert.Subject.Organization[0],
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
	}
}
