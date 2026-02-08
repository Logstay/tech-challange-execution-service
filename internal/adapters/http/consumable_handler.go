package http

import (
	"net/http"
	"strconv"

	cports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/errors"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	"github.com/gin-gonic/gin"
)

// SaveConsumable godoc
//
//	@Summary		Criar consumível
//	@Description	Cria um novo consumível no estoque
//	@Tags			Consumíveis
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SaveConsumableRequest	true	"Dados do consumível"
//	@Success		201		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/consumables [post]
func (h *Handler) SaveConsumable(c *gin.Context) {
	var req SaveConsumableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.SaveConsumable(c.Request.Context(), cports.SaveConsumable{
		Name:        null.StringFrom(req.Name),
		Description: null.StringFrom(req.Description),
		Price:       null.FloatFrom(req.Price),
		Code:        null.StringFrom(req.Code),
		Quantity:    null.IntFrom(int64(req.Quantity)),
	})

	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Consumable created successfully"})
}

// GetConsumable godoc
//
//	@Summary		Buscar consumível por UID
//	@Description	Retorna um consumível pelo seu identificador único
//	@Tags			Consumíveis
//	@Produce		json
//	@Param			uid	path		string	true	"UID do consumível"
//	@Success		200	{object}	GetConsumableResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/consumables/{uid} [get]
func (h *Handler) GetConsumable(c *gin.Context) {
	uid := c.Param("uid")

	result, err := h.services.GetConsumableByUID(c.Request.Context(), uid)
	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	resp := GetConsumableResponse{
		UID:         result.UID.String,
		Name:        result.Name.String,
		Description: result.Description.String,
		Code:        result.Code.String,
		Quantity:    result.Quantity.Int64,
		Price:       result.Price.Float64,
		CreatedAt:   result.CreatedAt.DateTime,
		UpdatedAt:   result.UpdatedAt.DateTime,
	}

	c.JSON(http.StatusOK, resp)
}

// ListConsumable godoc
//
//	@Summary		Listar consumíveis
//	@Description	Retorna uma lista paginada de consumíveis
//	@Tags			Consumíveis
//	@Produce		json
//	@Param			page	query		int		false	"Número da página"	default(1)
//	@Param			limit	query		int		false	"Itens por página"	default(10)
//	@Param			name	query		string	false	"Filtrar por nome"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/consumables [get]
func (h *Handler) ListConsumable(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	name := c.Query("name")

	result, total, err := h.services.ListConsumable(c.Request.Context(), cports.Filter{
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

// UpdateConsumable godoc
//
//	@Summary		Atualizar consumível
//	@Description	Atualiza os dados de um consumível existente
//	@Tags			Consumíveis
//	@Accept			json
//	@Produce		json
//	@Param			uid		path		string					true	"UID do consumível"
//	@Param			request	body		UpdateConsumableRequest	true	"Dados para atualização"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/consumables/{uid} [put]
func (h *Handler) UpdateConsumable(c *gin.Context) {
	uid := c.Param("uid")

	var req UpdateConsumableRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.services.UpdateConsumable(c.Request.Context(), cports.UpdateConsumable{
		UID:         null.StringFrom(uid),
		Name:        null.StringFrom(req.Name),
		Description: null.StringFrom(req.Description),
		Price:       null.FloatFrom(req.Price),
	})

	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consumable updated successfully"})
}

// DeleteConsumable godoc
//
//	@Summary		Deletar consumível
//	@Description	Remove um consumível pelo seu identificador único
//	@Tags			Consumíveis
//	@Produce		json
//	@Param			uid	path		string	true	"UID do consumível"
//	@Success		200	{object}	MessageResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/consumables/{uid} [delete]
func (h *Handler) DeleteConsumable(c *gin.Context) {
	uid := c.Param("uid")

	err := h.services.DeleteConsumable(c.Request.Context(), uid)
	if err != nil {
		errors.ErrorHandler(err, c)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consumable deleted successfully"})
}
