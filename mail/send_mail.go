package mail

import (
	"crypto/tls"
	"strconv"

	"net/mail"
	"net/smtp"

	"sentinel/config"
	"sentinel/logger"
	"sentinel/models"
	"sentinel/pkg/parseHtml"

	"github.com/xuri/excelize/v2"
	"gopkg.in/gomail.v2"
)

// SendMail is a function that sends an email with the given parameters
func SendMail(content *models.Mail, logs []models.Log, attachment *excelize.File) error {
	// Initialize mail server settings
	mailConfig := &models.MailConfig{
		Host:     config.C.Mail.Host,
		Port:     config.C.Mail.Port,
		Username: config.C.Mail.Username,
		Password: config.C.Mail.Password,
		FromName: config.C.Mail.FromName,
		FromMail: config.C.Mail.FromMail,
	}

	// Send mail with SMTP
	// err := SendMailWithSmtp(mailConfig, content, logs, attachment)
	// if err != nil {
	// 	logger.ERROR.Println("ERROR: ", err)
	// 	return err
	// }

	e := SendMailWithGomail(mailConfig, content, logs, attachment)
	if e != nil {
		logger.CLogger.Error("ERROR: ", e)
		return e
	}

	// If everything goes well, return nil
	logger.CLogger.Success("Email sent successfully")
	return nil
}

func SendMailWithSmtp(mailConfig *models.MailConfig, content *models.Mail, logs []models.Log, attach *excelize.File) error {
	// Set up authentication information.
	auth := smtp.PlainAuth(
		"",
		mailConfig.Username,
		mailConfig.Password,
		mailConfig.Host,
	)

	from := mail.Address{
		Name:    mailConfig.FromName,
		Address: mailConfig.FromMail,
	}

	body := parseHtml.LogTemplate(content, logs, "smtp")

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	err := smtp.SendMail(mailConfig.Host+":"+mailConfig.Port, auth, from.Address, []string(content.To), []byte(body))
	if err != nil {
		logger.CLogger.Error("ERROR: ", err)
		return err
	}

	return nil
}

// Send Mail with Gomail
func SendMailWithGomail(mailConfig *models.MailConfig, content *models.Mail, logs []models.Log, attach *excelize.File) error {
	m := gomail.NewMessage()

	m.SetHeader("From", mailConfig.FromMail)
	m.SetHeader("To", content.To...)
	m.SetHeader("Cc", content.Cc...)
	m.SetHeader("Subject", content.Subject)
	m.SetBody("text/html", parseHtml.LogTemplate(content, logs, "goMail"))

	port, err := strconv.Atoi(mailConfig.Port)
	if err != nil {
		return err
	}

	conn := gomail.NewDialer(mailConfig.Host, port, mailConfig.Username, mailConfig.Password)
	conn.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// ping, err := conn.Dial()
	// if err != nil {
	// 	fmt.Println("Sunucuya ping atılamadı ", err)
	// 	return err
	// }
	// defer ping.Close()

	if err := conn.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
