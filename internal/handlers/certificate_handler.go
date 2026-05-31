package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sentinel/config"
	"sentinel/internal/models"
	"sentinel/pkg/logger"
	"sentinel/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ListCertificates godoc
// ListCertificates This endpoint retrieves which domains been checked when /certificates/all endpoint called.
// @Summary Display domain list
// @Schemes
// @Description This endpoint retrieves the list of domains that have been checked when the /certificates/all endpoint was called.
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "List of domains successfully retrieved"
// @Failure 400 {object} models.APIResponse "List certificates failed due to invalid request"
// @Failure 422 {object} models.APIResponse "List certificates failed due to invalid request"
// @Failure 500 {object} models.APIResponse "List certificates failed due to internal server error"
// @Router /domains [GET]
func (ss *Sentinel) ListCertificates(c *gin.Context) (int, interface{}, error) {
	// Get domain list from DataService (DB or static data based on config)
	domainList := utils.DomainList
	if utils.GlobalDataService != nil {
		domainList = utils.GlobalDataService.GetDomainList()
	}

	return http.StatusOK, domainList, nil
}

// GetAllExpirations godoc
// GetAllExpirations handles the request to check all domain certificates
// @Summary Check all domain SSL/TLS certificates
// @Schemes
// @Description This endpoint checks all domain SSL/TLS certificates and sends an email with the results.
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "All certificates checked successfully"
// @Failure 400 {object} models.APIResponse "Certificate check failed due to invalid request body"
// @Failure 422 {object} models.APIResponse "Certificate check failed due to invalid request body"
// @Failure 500 {object} models.APIResponse "Certificate check failed due to internal server error"
// @Router /certificates/scan [GET]
func (ss *Sentinel) GetAllExpirations(c *gin.Context) (int, interface{}, error) {
	logger.CLogger.Info("INFO: Checking all certificates concurrently...")

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	// Get domain list from DataService (DB or static data based on config)
	domainList := utils.DomainList
	if utils.GlobalDataService != nil {
		domainList = utils.GlobalDataService.GetDomainList()
	}

	if len(domainList) == 0 {
		logger.CLogger.Warn("WARN: No domains found to check")
		return http.StatusOK, nil, fmt.Errorf("no domains configured to check")
	}

	logger.CLogger.Infof("INFO: Checking %d domains", len(domainList))

	// Use concurrent certificate checker with worker pool
	checker := utils.NewCertificateChecker(10) // 10 concurrent workers
	logs := checker.CheckCertificates(ctx, domainList, config.C.App.Expire)

	// Send to Mail if there are logs
	if len(logs) > 0 {
		f := utils.SetChangesToExcel(logs)
		if f == nil {
			logger.CLogger.Error("ERROR: Failed to create Excel file")
			return http.StatusInternalServerError, nil, fmt.Errorf("failed to create Excel file")
		}
		logger.CLogger.Info("INFO: Excel file created successfully")
		utils.SendMailWithAttachment(logs, f)
	} else {
		logger.CLogger.Info("INFO: No expiring certificates found")
	}

	// Return all logs
	return http.StatusOK, logs, nil
}

// GetCertificateInfo godoc
// GetCertificateInfo handles the given domain certificate check request
// @Summary Check a given domain certificate
// @Schemes
// @Description This endpoint checks the SSL/TLS certificate of a given domain and returns its details. Only one domain can be checked at a time.
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param domain query string true "Domain"
// @Param expire query int false "Expire days, default is config.C.App.Expire value"
// @Success 200 {object} models.APIResponse "Certificate details retrieved successfully"
// @Failure 400 {object} models.APIResponse "Certificate check failed due to invalid request"
// @Failure 422 {object} models.APIResponse "Certificate check failed due to invalid request"
// @Failure 500 {object} models.APIResponse "Certificate check failed due to internal server error"
// @Router /certificates/check [get]
func (ss *Sentinel) GetCertificateInfo(c *gin.Context) (int, interface{}, error) {
	// Get domain from the request parameters
	domain := c.Query("domain")
	if domain == "" {
		return http.StatusBadRequest, nil, fmt.Errorf("domain is required")
	}

	// Get expire days from the query parameters, default is 30 days
	expireDays := config.C.App.Expire
	if expireQuery := c.Query("expire"); expireQuery != "" {
		if days, err := strconv.Atoi(expireQuery); err == nil && days > 0 {
			expireDays = days
		} else {
			return http.StatusBadRequest, nil, fmt.Errorf("invalid expire days, must be a positive integer")
		}
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	var logs []models.Log
	var isExpired bool
	var data *models.Log

	// Use retry logic with context
	retryConfig := utils.DefaultRetryConfig()
	err := utils.RetryWithContext(ctx, retryConfig, func(attemptCtx context.Context) error {
		isExpired, data = utils.CheckDomainCertificateWithContext(attemptCtx, domain, expireDays)
		if !isExpired && data == nil {
			return fmt.Errorf("failed to retrieve certificate")
		}
		return nil
	}, domain)

	if err != nil {
		return http.StatusOK, nil, fmt.Errorf("%s domain expires in more than %d days or check failed", domain, expireDays)
	}

	if isExpired && data != nil {
		logs = append(logs, *data)
	} else if data != nil {
		logger.CLogger.Infof("INFO: %s - %s", domain, data.Message)
		return http.StatusOK, nil, fmt.Errorf("%s domain expires in more than %d days", domain, expireDays)
	}

	return http.StatusOK, logs, nil
}
