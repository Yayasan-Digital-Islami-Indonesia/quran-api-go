package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"quran-api-go/internal/domain/healthcheck"
	"quran-api-go/pkg/response"
)

type HealthCheckHandler struct {
	service healthcheck.HealthCheckService
}

func NewHealthCheckHandler(service healthcheck.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{service: service}
}

// HealthCheck godoc
// @Summary     Health check
// @Description Check if the API is running
// @Tags        Health
// @Produce     json
// @Success     200  {object}  response.SuccessResponse{data=healthcheck.HealthCheck}
// @Failure     503  {object}  response.ErrorResponse
// @Router      /health [get]
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	health, err := h.service.HealthCheck(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	response.Success(c, health)
}

// ReadyCheck godoc
// @Summary     Readiness check
// @Description Check if the API is ready to serve requests (database connectivity)
// @Tags        Health
// @Produce     json
// @Success     200  {object}  response.SuccessResponse{data=healthcheck.HealthCheck}
// @Failure     503  {object}  response.ErrorResponse
// @Router      /health/ready [get]
func (h *HealthCheckHandler) ReadyCheck(c *gin.Context) {
	ready, err := h.service.ReadyCheck(c.Request.Context())
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	response.Success(c, ready)
}
