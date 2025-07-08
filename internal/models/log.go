package models

import "time"

// SSL Certificate Model
type Log struct {
	// ID                 int       `json:"id" gorm:"primaryKey"`
	Version            int       `json:"version" gorm:"version"`
	SerialNumber       string    `json:"serial_number" gorm:"serial_number"`
	Subject            string    `json:"subject" gorm:"subject"`
	IssuerSubject      string    `json:"issuer_subject" gorm:"issuer_subject"`
	Domain             string    `json:"domain" gorm:"domain"`
	Port               int       `json:"port" gorm:"port"`
	CommonName         string    `json:"common_name" gorm:"common_name"`
	Organization       string    `json:"organization" gorm:"organization"`
	IssuedOn           time.Time `json:"issued_on" gorm:"issued_on"`
	ExpiresOn          time.Time `json:"expires_on" gorm:"expires_on"`
	CertificateData    string    `json:"certificate_data" gorm:"certificate_data"`
	SignatureAlgorithm string    `json:"signature_algorithm" gorm:"signature_algorithm"`
	SubjectKeyID       string    `json:"subject_key_id" gorm:"subject_key_id"`
	AuthorityKeyID     string    `json:"authority_key_id" gorm:"authority_key_id"`
	IsCA               bool      `json:"is_ca" gorm:"is_ca"`
	Issuer             string    `json:"issuer" gorm:"issuer"`
	IsExpired          bool      `json:"is_expired" gorm:"is_expired"`
	Message            string    `json:"message" gorm:"message"`
	Status             int       `json:"status" gorm:"status"` // 0: Not Expired, 1: Expired 2: Time Out
}
