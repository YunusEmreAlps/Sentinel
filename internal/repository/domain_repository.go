package repository

import (
	"sentinel/internal/models"

	"gorm.io/gorm"
)

// DomainRepository handles domain database operations
type DomainRepository struct {
	db *gorm.DB
}

// NewDomainRepository creates a new domain repository
func NewDomainRepository(db *gorm.DB) *DomainRepository {
	return &DomainRepository{db: db}
}

// GetAll returns all active domains
func (r *DomainRepository) GetAll() ([]models.Domain, error) {
	var domains []models.Domain
	err := r.db.Where("active = ?", true).Find(&domains).Error
	return domains, err
}

// GetActiveDomainStrings returns all active domain strings
func (r *DomainRepository) GetActiveDomainStrings() ([]string, error) {
	var domains []models.Domain
	if err := r.db.Where("active = ?", true).Find(&domains).Error; err != nil {
		return nil, err
	}

	domainStrings := make([]string, 0, len(domains))
	for _, domain := range domains {
		domainStrings = append(domainStrings, domain.Domain)
	}

	return domainStrings, nil
}

// GetByID returns a domain by ID
func (r *DomainRepository) GetByID(id uint) (*models.Domain, error) {
	var domain models.Domain
	err := r.db.First(&domain, id).Error
	return &domain, err
}

// Create creates a new domain
func (r *DomainRepository) Create(domain *models.Domain) error {
	return r.db.Create(domain).Error
}

// Update updates an existing domain
func (r *DomainRepository) Update(domain *models.Domain) error {
	return r.db.Save(domain).Error
}

// Delete soft deletes a domain by setting active to false
func (r *DomainRepository) Delete(id uint) error {
	return r.db.Model(&models.Domain{}).Where("id = ?", id).Update("active", false).Error
}

// HardDelete permanently deletes a domain
func (r *DomainRepository) HardDelete(id uint) error {
	return r.db.Delete(&models.Domain{}, id).Error
}
