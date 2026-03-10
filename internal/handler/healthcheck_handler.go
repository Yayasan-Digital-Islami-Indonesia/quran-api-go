package handler

import (
	"net/http"
	"quran-api-go/internal/domain/healthcheck"

	"github.com/gin-gonic/gin"
)

type HealthCheckHandler struct {
	service healthcheck.HealthCheckService
}

func NewHealthCheckHandler(service healthcheck.HealthCheckService) *HealthCheckHandler {
	return &HealthCheckHandler{service: service}
}

func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	health, err := h.service.HealthCheck(c)
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	c.JSON(http.StatusOK, health)
}

func (h *HealthCheckHandler) ReadyCheck(c *gin.Context) {
	ready, err := h.service.ReadyCheck(c)
	if err != nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}
	c.JSON(http.StatusOK, ready)
}
