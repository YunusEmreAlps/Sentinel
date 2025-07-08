package handlers

import (
	"net/http"
	"sentinel/pkg/utils"

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
// @Success 200 {object} RespondJson "List of domains successfully retrieved"
// @Failure 400 {object} RespondJson "List certificates failed due to invalid request"
// @Failure 422 {object} RespondJson "List certificates failed due to invalid request"
// @Failure 500 {object} RespondJson "List certificates failed due to internal server error"
// @Router /certificates [get]
func (ss *Sentinel) ListCertificates(c *gin.Context) (int, interface{}, error) {
	return http.StatusOK, utils.DomainList, nil
}
