package models

import (
	"time"
)

// Domain represents a monitored domain in the database
type Domain struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Domain    string    `gorm:"unique;not null" json:"domain" binding:"required"` // Supports all formats: google.com, https://example.com, domain.com:443
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Domain model
func (Domain) TableName() string {
	return "domains"
}
