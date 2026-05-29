package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"sentinel/config"
	"sentinel/pkg/logger"
	"sentinel/pkg/utils"

	"github.com/gin-gonic/gin"
)

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

	// Use concurrent certificate checker with worker pool
	checker := utils.NewCertificateChecker(10) // 10 concurrent workers
	logs := checker.CheckCertificates(ctx, utils.DomainList, config.C.App.Expire)

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
