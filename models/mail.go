package models

type MailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	FromMail string
	FromName string
}

type Mail struct {
	Sender  string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
}

type MailContent struct {
	Nonce       string `json:"nonce" bson:"nonce" validate:"required" binding:"required" example:"123"`
	To          string `json:"to" bson:"to" validate:"required,email" binding:"required,email" example:"test@domain.com"`
	ToName      string `json:"to_name" bson:"to_name" validate:"required" binding:"required" example:"testuser"`
	Subject     string `json:"subject" bson:"subject" validate:"required" binding:"required" example:"Test mail"`
	Body        string `json:"body" bson:"body" validate:"required" binding:"required" example:"This is a test mail"`
	MimeVersion string `json:"mime_version" bson:"mime_version" validate:"required" binding:"required" example:"1.0"`
	ContentType string `json:"content_type" bson:"content_type" validate:"required" binding:"required" example:"text/html"`
}
