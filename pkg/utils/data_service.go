package utils

import (
	"sentinel/config"
	"sentinel/internal/repository"
	"sentinel/pkg/logger"

	"gorm.io/gorm"
)

// DataService handles data retrieval from DB or static data
type DataService struct {
	domainRepo *repository.DomainRepository
	dbActive   bool
}

// NewDataService creates a new data service
func NewDataService(db *gorm.DB) *DataService {
	return &DataService{
		domainRepo: repository.NewDomainRepository(db),
		dbActive:   config.C.DB.Active && db != nil,
	}
}

// GetDomainList returns domain list from DB if active, otherwise from static data
func (s *DataService) GetDomainList() []string {
	if s.dbActive {
		domains, err := s.domainRepo.GetActiveDomainStrings()
		if err != nil {
			logger.CLogger.Errorf("Failed to get domains from DB, falling back to static data: %v", err)
			return DomainList
		}

		if len(domains) == 0 {
			logger.CLogger.Warn("No active domains found in DB, using static data")
			return DomainList
		}

		logger.CLogger.Infof("Loaded %d domains from database", len(domains))
		return domains
	}

	logger.CLogger.Info("Using static domain list from data.go")
	return DomainList
}

// GetDomainRepo returns the domain repository (only if DB is active)
func (s *DataService) GetDomainRepo() *repository.DomainRepository {
	return s.domainRepo
}

// Global data service instance
var GlobalDataService *DataService
