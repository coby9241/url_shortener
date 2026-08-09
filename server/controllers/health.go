package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthController provides a simple health check endpoint.
type HealthController struct{}

// NewHealthController creates a new HealthController.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck returns a 200 OK response with a status message.
func (hc *HealthController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}