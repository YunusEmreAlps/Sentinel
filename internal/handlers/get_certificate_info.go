package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"sentinel/config"
	"sentinel/internal/models"
	"sentinel/pkg/logger"
	"sentinel/pkg/utils"

	"github.com/gin-gonic/gin"
)

// GetCertificateInfo godoc
// GetCertificateInfo handles the given domain certificate check request
// @Summary Check a given domain certificate
// @Schemes
// @Description This endpoint checks the SSL/TLS certificate of a given domain and returns its details. Only one domain can be checked at a time.
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param domain path string true "yunusemrealpu.netlify.app:443"
// @Param expire query int false "Expire days, default is config.C.App.Expire value"
// @Success 200 {object} models.APIResponse "Certificate details retrieved successfully"
// @Failure 400 {object} models.APIResponse "Certificate check failed due to invalid request"
// @Failure 422 {object} models.APIResponse "Certificate check failed due to invalid request"
// @Failure 500 {object} models.APIResponse "Certificate check failed due to internal server error"
// @Router /certificates/{domain} [get]
func (ss *Sentinel) GetCertificateInfo(c *gin.Context) (int, interface{}, error) {
	// Get domain from the request parameters
	domain := c.Param("domain")
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
