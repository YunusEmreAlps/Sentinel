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

// GetAllExpirations godoc
// GetAllExpirations handles the request to check all domain certificates
// @Summary Check all domain SSL/TLS certificates
// @Schemes
// @Description This endpoint checks all domain SSL/TLS certificates and sends an email with the results.
// @Tags Certificates
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} RespondJson "All certificates checked successfully"
// @Failure 400 {object} RespondJson "Certificate check failed due to invalid request body"
// @Failure 422 {object} RespondJson "Certificate check failed due to invalid request body"
// @Failure 500 {object} RespondJson "Certificate check failed due to internal server error"
// @Router /certificates/all [get]
func (ss *Sentinel) GetAllExpirations(c *gin.Context) (int, interface{}, error) {
	// Step 1: Get all certificates from the utility/data.go
	var logs []models.Log

	fmt.Println("INFO: Checking all certificates...")

	for _, domain := range utils.DomainList {

		const maxRetries = 3
		for i := 0; i < maxRetries; i++ {
			isOK, data := utils.CheckDomainCertificate(domain, config.C.App.Expire)
			if isOK && data != nil {
				// Successfully retrieved certificate, break out of the loop
				logs = append(logs, *data)
				break
			} else {
				if data != nil {
					// Certificate will not expired in 30 days
					logger.CLogger.Info("INFO: ", domain+" - "+data.Message)
					break
				} else {
					// Connection Error
					logger.CLogger.Error("ERROR: ", domain+" - Connection Error Attempt: "+strconv.Itoa(i+1)+"/"+strconv.Itoa(maxRetries))
				}
			}
			// Wait for a brief period before retrying
			time.Sleep(2 * time.Second)
		}
	}

	// Send to Mail
	f := utils.SetChangesToExcel(logs)
	if f == nil {
		logger.CLogger.Error("ERROR: ", "Failed to create Excel file")
		return http.StatusInternalServerError, nil, fmt.Errorf("failed to create Excel file")
	} else {
		logger.CLogger.Info("INFO: ", "Excel file created successfully")
		utils.SendMailWithAttachment(logs, f)
	}

	// Step 2: Return all logs
	return http.StatusOK, logs, nil
}
