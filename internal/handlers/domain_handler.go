package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"sentinel/config"
	"sentinel/internal/models"
	"sentinel/pkg/logger"
	"sentinel/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetAllDomains godoc
// @Summary Get all domains from database
// @Description Returns all active domains from database (only works if DB is active)
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "List of domains successfully retrieved"
// @Failure 403 {object} models.APIResponse "Database is not active"
// @Failure 500 {object} models.APIResponse "Failed to retrieve domains"
// @Router /domains/manage [get]
func (ss *Sentinel) GetAllDomains(c *gin.Context) (int, interface{}, error) {
	if !config.C.DB.Active || utils.GlobalDataService == nil {
		return http.StatusForbidden, nil, fmt.Errorf("database feature is not active")
	}

	domains, err := utils.GlobalDataService.GetDomainRepo().GetAll()
	if err != nil {
		logger.CLogger.Errorf("Failed to get domains: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to retrieve domains")
	}

	return http.StatusOK, domains, nil
}

// GetDomainByID godoc
// @Summary Get domain by ID
// @Description Returns a specific domain by ID (only works if DB is active)
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Domain ID"
// @Success 200 {object} models.APIResponse "Domain successfully retrieved"
// @Failure 400 {object} models.APIResponse "Invalid domain ID"
// @Failure 403 {object} models.APIResponse "Database is not active"
// @Failure 404 {object} models.APIResponse "Domain not found"
// @Failure 500 {object} models.APIResponse "Failed to retrieve domain"
// @Router /domains/manage/{id} [get]
func (ss *Sentinel) GetDomainByID(c *gin.Context) (int, interface{}, error) {
	if !config.C.DB.Active || utils.GlobalDataService == nil {
		return http.StatusForbidden, nil, fmt.Errorf("database feature is not active")
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid domain ID")
	}

	domain, err := utils.GlobalDataService.GetDomainRepo().GetByID(uint(id))
	if err != nil {
		logger.CLogger.Errorf("Failed to get domain %d: %v", id, err)
		return http.StatusNotFound, nil, fmt.Errorf("domain not found")
	}

	return http.StatusOK, domain, nil
}

// CreateDomain godoc
// @Summary Create a new domain
// @Description Add a new domain to monitor (only works if DB is active). Supports formats: google.com, https://example.com, domain.com:443
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param domain body models.Domain true "Domain object"
// @Success 201 {object} models.APIResponse "Domain successfully created"
// @Failure 400 {object} models.APIResponse "Invalid request body"
// @Failure 403 {object} models.APIResponse "Database is not active"
// @Failure 500 {object} models.APIResponse "Failed to create domain"
// @Router /domains/manage [post]
func (ss *Sentinel) CreateDomain(c *gin.Context) (int, interface{}, error) {
	if !config.C.DB.Active || utils.GlobalDataService == nil {
		return http.StatusForbidden, nil, fmt.Errorf("database feature is not active")
	}

	var domain models.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid request body: %v", err)
	}

	// Validate domain format
	if _, err := utils.NormalizeDomain(domain.Domain); err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid domain format: %v", err)
	}

	if err := utils.GlobalDataService.GetDomainRepo().Create(&domain); err != nil {
		logger.CLogger.Errorf("Failed to create domain: %v", err)
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to create domain")
	}

	logger.CLogger.Infof("Domain created successfully: %s", domain.Domain)
	return http.StatusCreated, domain, nil
}

// UpdateDomain godoc
// @Summary Update a domain
// @Description Update an existing domain (only works if DB is active)
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Domain ID"
// @Param domain body models.Domain true "Domain object"
// @Success 200 {object} models.APIResponse "Domain successfully updated"
// @Failure 400 {object} models.APIResponse "Invalid request"
// @Failure 403 {object} models.APIResponse "Database is not active"
// @Failure 404 {object} models.APIResponse "Domain not found"
// @Failure 500 {object} models.APIResponse "Failed to update domain"
// @Router /domains/manage/{id} [put]
func (ss *Sentinel) UpdateDomain(c *gin.Context) (int, interface{}, error) {
	if !config.C.DB.Active || utils.GlobalDataService == nil {
		return http.StatusForbidden, nil, fmt.Errorf("database feature is not active")
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid domain ID")
	}

	var domain models.Domain
	if err := c.ShouldBindJSON(&domain); err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid request body: %v", err)
	}

	// Validate domain format
	if _, err := utils.NormalizeDomain(domain.Domain); err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid domain format: %v", err)
	}

	domain.ID = uint(id)
	if err := utils.GlobalDataService.GetDomainRepo().Update(&domain); err != nil {
		logger.CLogger.Errorf("Failed to update domain %d: %v", id, err)
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to update domain")
	}

	logger.CLogger.Infof("Domain updated successfully: %s", domain.Domain)
	return http.StatusOK, domain, nil
}

// DeleteDomain godoc
// @Summary Delete a domain (soft delete)
// @Description Soft delete a domain by setting active to false (only works if DB is active)
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Domain ID"
// @Success 200 {object} models.APIResponse "Domain successfully deleted"
// @Failure 400 {object} models.APIResponse "Invalid domain ID"
// @Failure 403 {object} models.APIResponse "Database is not active"
// @Failure 500 {object} models.APIResponse "Failed to delete domain"
// @Router /domains/manage/{id} [delete]
func (ss *Sentinel) DeleteDomain(c *gin.Context) (int, interface{}, error) {
	if !config.C.DB.Active || utils.GlobalDataService == nil {
		return http.StatusForbidden, nil, fmt.Errorf("database feature is not active")
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid domain ID")
	}

	if err := utils.GlobalDataService.GetDomainRepo().Delete(uint(id)); err != nil {
		logger.CLogger.Errorf("Failed to delete domain %d: %v", id, err)
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to delete domain")
	}

	logger.CLogger.Infof("Domain deleted successfully: ID %d", id)
	return http.StatusOK, gin.H{"message": "Domain deleted successfully"}, nil
}

// HardDeleteDomain godoc
// @Summary Hard delete a domain
// @Description Permanently delete a domain from database (only works if DB is active)
// @Tags Domain Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Domain ID"
// @Success 200 {object} models.APIResponse "Domain permanently deleted"
// @Failure 400 {object} models.APIResponse "Invalid domain ID"
// @Failure 403 {object} models.APIResponse "Database is not active"
// @Failure 500 {object} models.APIResponse "Failed to delete domain"
// @Router /domains/manage/{id}/hard [delete]
func (ss *Sentinel) HardDeleteDomain(c *gin.Context) (int, interface{}, error) {
	if !config.C.DB.Active || utils.GlobalDataService == nil {
		return http.StatusForbidden, nil, fmt.Errorf("database feature is not active")
	}

	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("invalid domain ID")
	}

	if err := utils.GlobalDataService.GetDomainRepo().HardDelete(uint(id)); err != nil {
		logger.CLogger.Errorf("Failed to hard delete domain %d: %v", id, err)
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to permanently delete domain")
	}

	logger.CLogger.Infof("Domain permanently deleted: ID %d", id)
	return http.StatusOK, gin.H{"message": "Domain permanently deleted"}, nil
}
