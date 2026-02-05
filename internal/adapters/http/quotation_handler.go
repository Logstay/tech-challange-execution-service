package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetQuotation godoc
//
//	@Summary		📄 Buscar Orçamento
//	@Description	Busca os detalhes de um orçamento existente pelo ID da preferência
//	@Tags			💰 Pagamentos
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"ID da Preferência (Orçamento)"
//	@Success		200	{object}	QuotationResponse	"Orçamento encontrado"
//	@Failure		404	{object}	httperr.RestErr		"Orçamento não encontrado"
//	@Failure		500	{object}	httperr.RestErr		"Erro interno do servidor"
//	@Router			/v1/quotations/{id} [get]
func (h *Handler) GetQuotation(c *gin.Context) {
	var (
		preferenceID = c.Param("id")
		total        float64
	)
	if preferenceID == "" {
		c.JSON(http.StatusBadRequest, "preference_id is required")
		return
	}

	ctx := c.Request.Context()
	preference, err := h.services.GetQuotationByID(ctx, preferenceID)
	if err != nil {
		c.JSON(http.StatusNotFound, "quotation not found")
		return
	}

	for _, item := range preference.Items {
		total += item.UnitPrice * float64(item.Quantity)
	}

	response := QuotationResponse{
		ID:                preference.ID,
		PaymentLink:       preference.InitPoint,
		SandboxLink:       preference.SandboxInitPoint,
		ExternalReference: preference.ExternalReference,
		TotalAmount:       total,
		Currency:          "BRL",
		Status:            preference.Status,
		CreatedAt:         preference.DateCreated,
		Items:             preference.Items,
		Metadata:          preference.Metadata,
	}

	if preference.ExpirationDateTo != nil {
		response.ExpiresAt = *preference.ExpirationDateTo
	}

	c.JSON(http.StatusOK, response)
}
