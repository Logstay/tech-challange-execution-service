package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
//
//	@Summary		🏥 Health Check
//	@Description	Verifica o status de saúde do serviço de cobrança
//	@Tags			🔧 Sistema
//	@Produce		json
//	@Success		200	{object}	map[string]string	"Serviço funcionando corretamente"
//	@Router			/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "tech-challange-wherehouse-service",
	})
}
