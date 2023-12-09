package parseHtml

import (
	"bytes"
	"fmt"
	"html/template"
	"sentinel/config"
	"sentinel/logger"
	"sentinel/models"
	"strings"
	"time"
)

func LogTemplate(content *models.Mail, logs []models.Log, l string) string {
	// Send HTML Email with Excel File as Attachment
	var templateBuffer bytes.Buffer

	if l == "smtp" {
		//Senders, receivers and subject
		templateBuffer.WriteString(fmt.Sprintf("From: %s\r\n", content.Sender))
		templateBuffer.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(content.To, ";")))
		templateBuffer.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(content.Cc, ";")))
		templateBuffer.WriteString(fmt.Sprintf("Bcc: %s\r\n", strings.Join(content.Bcc, ";")))
		templateBuffer.WriteString(fmt.Sprintf("Subject: %s\r\n", content.Subject))
	}

	t, err := template.ParseFiles("./templates/log.html")
	if err != nil {
		logger.CLogger.Error("ERROR: ", err)
		return ""
	}

	// Multipart MIME Header
	if l == "smtp" {
		mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		templateBuffer.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", "Kodeks Error Logs", mimeHeaders)))
	}

	// Execute the template
	r := t.Execute(&templateBuffer, struct {
		AppName     string
		CurrentTime string
		Logs        []models.Log
	}{
		AppName:     config.C.App.TargetApp,
		CurrentTime: time.Now().Format("15:04:05 02.01.2006"),
		Logs:        logs,
	})
	if r != nil {
		logger.CLogger.Error("ERROR: ", r)
		return ""
	}

	return templateBuffer.String()
}
