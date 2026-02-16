package http

import (
	"net/http"
	"strconv"

	sdports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/gin-gonic/gin"
)

// SaveServiceDefinition godoc
//
//	@Summary		Criar definição de serviço
//	@Description	Cria uma nova definição de serviço
//	@Tags			Definições de Serviço
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SaveServiceDefinitionRequest	true	"Dados da definição de serviço"
//	@Success		201		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/service-definitions [post]
func (h *Handler) SaveServiceDefinition(c *gin.Context) {
	var req SaveServiceDefinitionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.SaveServiceDefinition(c.Request.Context(), sdports.SaveServiceDefinition{
		Name:        null.StringFrom(req.Name),
		Description: null.StringFrom(req.Description),
		Price:       null.FloatFrom(req.Price),
	})

	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Service definition created successfully"})
}

// GetServiceDefinition godoc
//
//	@Summary		Buscar definição de serviço por UID
//	@Description	Retorna uma definição de serviço pelo seu identificador único
//	@Tags			Definições de Serviço
//	@Produce		json
//	@Param			uid	path		string	true	"UID da definição de serviço"
//	@Success		200	{object}	GetServiceDefinitionResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/service-definitions/{uid} [get]
func (h *Handler) GetServiceDefinition(c *gin.Context) {
	uid := c.Param("uid")

	result, err := h.services.GetServiceDefinitionByUID(c.Request.Context(), uid)
	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	resp := GetServiceDefinitionResponse{
		UID:         result.UID.String,
		Name:        result.Name.String,
		Description: result.Description.String,
		Price:       result.Price.Float64,
		CreatedAt:   result.CreatedAt.DateTime,
		UpdatedAt:   result.UpdatedAt.DateTime,
	}

	c.JSON(http.StatusOK, resp)
}

// ListServiceDefinition godoc
//
//	@Summary		Listar definições de serviço
//	@Description	Retorna uma lista paginada de definições de serviço
//	@Tags			Definições de Serviço
//	@Produce		json
//	@Param			page	query		int		false	"Número da página"	default(1)
//	@Param			limit	query		int		false	"Itens por página"	default(10)
//	@Param			name	query		string	false	"Filtrar por nome"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/service-definitions [get]
func (h *Handler) ListServiceDefinition(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	name := c.Query("name")

	result, total, err := h.services.ListServiceDefinition(c.Request.Context(), sdports.Filter{
		Page:  null.IntFrom(page),
		Limit: null.IntFrom(limit),
		Name:  null.StringFrom(name),
	})

	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateServiceDefinition godoc
//
//	@Summary		Atualizar definição de serviço
//	@Description	Atualiza os dados de uma definição de serviço existente
//	@Tags			Definições de Serviço
//	@Accept			json
//	@Produce		json
//	@Param			uid		path		string							true	"UID da definição de serviço"
//	@Param			request	body		UpdateServiceDefinitionRequest	true	"Dados para atualização"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/service-definitions/{uid} [put]
func (h *Handler) UpdateServiceDefinition(c *gin.Context) {
	uid := c.Param("uid")

	var req UpdateServiceDefinitionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.UpdateServiceDefinition(c.Request.Context(), sdports.UpdateServiceDefinition{
		UID:         null.StringFrom(uid),
		Name:        null.StringFrom(req.Name),
		Description: null.StringFrom(req.Description),
		Price:       null.FloatFrom(req.Price),
	})

	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service definition updated successfully"})
}

// DeleteServiceDefinition godoc
//
//	@Summary		Deletar definição de serviço
//	@Description	Remove uma definição de serviço pelo seu identificador único
//	@Tags			Definições de Serviço
//	@Produce		json
//	@Param			uid	path		string	true	"UID da definição de serviço"
//	@Success		200	{object}	MessageResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/service-definitions/{uid} [delete]
func (h *Handler) DeleteServiceDefinition(c *gin.Context) {
	uid := c.Param("uid")

	err := h.services.DeleteServiceDefinition(c.Request.Context(), uid)
	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service definition deleted successfully"})
}
