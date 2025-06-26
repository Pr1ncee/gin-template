package handlers

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HealthCheckHandler struct {
	logger *log.Logger
}

func NewHealthCheckHandler(logger *log.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{logger: logger}
}

// HealthCheck godoc
// @Summary A method which indicates whether the application is healthy
// @Description Returns a response indicating a healthy application
// @Tags health-check
// @Produce json
// @Success 200 {object} HealthCheckResponse
// @Router /health-check [get]
func (h *HealthCheckHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"Status": "success"})
}

type HealthCheckResponse struct {
	Status string `json:"Status" example:"success"`
}
