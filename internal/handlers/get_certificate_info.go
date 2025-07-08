package handlers

import (
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
// @Success 200 {object} RespondJson "Certificate details retrieved successfully"
// @Failure 400 {object} RespondJson "Certificate check failed due to invalid request"
// @Failure 422 {object} RespondJson "Certificate check failed due to invalid request"
// @Failure 500 {object} RespondJson "Certificate check failed due to internal server error"
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

	var logs []models.Log
	const maxRetries = 3
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		isOK, data := utils.CheckDomainCertificate(domain, expireDays)
		if isOK && data != nil {
			// Successfully retrieved certificate
			logs = append(logs, *data)
			break
		}

		if data != nil {
			// Certificate is valid and not expiring soon
			logger.CLogger.Info(fmt.Sprintf("INFO: %s - %s", domain, data.Message))
			break
		}

		// Log connection error and retry
		logger.CLogger.Error(fmt.Sprintf("ERROR: %s - Connection Error Attempt: %d/%d", domain, i+1, maxRetries))
		time.Sleep(retryDelay)
	}

	// If no logs were added, return an error
	if len(logs) == 0 {
		return http.StatusOK, nil, fmt.Errorf("%s domain expires in more than %d days", domain, expireDays)
	}

	return http.StatusOK, logs, nil
}
