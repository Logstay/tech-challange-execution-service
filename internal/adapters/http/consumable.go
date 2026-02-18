package http

import (
	"net/http"

	cports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SaveConsumable godoc
//
//	@Security		BearerAuth
//	@Summary		Criar consumível
//	@Description	Cria um novo consumível no estoque
//	@Tags			Consumíveis
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SaveConsumableRequest	true	"Dados do consumível"
//	@Success		201		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/consumables [post]
func (h *Handler) SaveConsumable(c *gin.Context) {

	var (
		parentCtx = tracing.NewRelicContext(c)
		in        = SaveConsumableRequest{}
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-save-consumable")
	defer func() {
		span.End()
	}()

	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	span.AddEvent("http-handler-save-consumable-params", tracing.Attributes{
		"model": in,
	})

	err := h.services.SaveConsumable(ctx, cports.SaveConsumable{
		Name:     null.StringFrom(in.Name),
		Price:    null.FloatFrom(in.Price),
		Quantity: null.IntFrom(in.Quantity),
		Description: null.StringFrom(
			in.Description,
		),
	})
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-save-consumable-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})

		return
	}

	span.AddEvent("http-handler-save-consumable-success", nil)

	c.JSON(http.StatusCreated, nil)

}

// GetConsumable godoc
//
//	@Security		BearerAuth
//	@Summary		Buscar consumível por UID
//	@Description	Retorna um consumível pelo seu identificador único
//	@Tags			Consumíveis
//	@Produce		json
//	@Param			uid	path		string	true	"UID do consumível"
//	@Success		200	{object}	GetConsumableResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/consumables/{uid} [get]
func (h *Handler) GetConsumable(c *gin.Context) {
	var (
		parentCtx = tracing.NewRelicContext(c)
		uid       = c.Param("uid")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-get-consumable")
	defer func() {
		span.End()
	}()

	if uid == "" {
		c.JSON(http.StatusBadRequest, "uid é obrigatorio")
		return
	}

	_, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, "uid invalido")
		return
	}

	resp, err := h.services.GetConsumableByUID(ctx, uid)
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-get-consumable-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})
		return
	}

	result := GetConsumableResponse{
		UID:         resp.UID.String,
		Name:        resp.Name.String,
		Price:       float64ToReais(resp.Price.Float64),
		Quantity:    resp.Quantity.Int64,
		Code:        resp.Code.String,
		Description: resp.Description.String,
		CreatedAt:   resp.CreatedAt.Ptr(),
		UpdatedAt:   resp.UpdatedAt.Ptr(),
	}

	span.AddEvent("http-handler-get-consumable-success", nil)

	c.JSON(http.StatusOK, result)

}

// ListConsumable godoc
//
//	@Security		BearerAuth
//	@Summary		Listar consumíveis
//	@Description	Retorna uma lista paginada de consumíveis
//	@Tags			Consumíveis
//	@Produce		json
//	@Param			page	query		int		false	"Número da página"	default(0)
//	@Param			limit	query		int		false	"Itens por página"	default(10)
//	@Param			name	query		string	false	"Filtrar por nome"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/consumables [get]
func (h *Handler) ListConsumable(c *gin.Context) {

	page, limit := pageFromRequest(c.Request)

	var (
		parentCtx = tracing.NewRelicContext(c)
		name      = c.Query("name")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-list-consumable")
	defer func() {
		span.End()
	}()

	resp, total, err := h.services.ListConsumable(ctx, cports.Filter{
		Name:  null.StringFrom(name),
		Page:  null.IntFrom(page),
		Limit: null.IntFrom(limit),
	})
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-list-consumable-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})
		return
	}

	result := ListConsumableResponse{}
	for _, r := range resp {
		result.Data = append(result.Data, GetConsumableResponse{
			UID:         r.UID.String,
			Name:        r.Name.String,
			Price:       float64ToReais(r.Price.Float64),
			Code:        r.Code.String,
			Description: r.Description.String,
			Quantity:    r.Quantity.Int64,
			CreatedAt:   r.CreatedAt.Ptr(),
			UpdatedAt:   r.UpdatedAt.Ptr(),
		})
	}

	result.Pagination = Pagination{
		Page:  page,
		Limit: limit,
		Total: total,
	}

	span.AddEvent("http-handler-list-consumable-success", nil)

	c.JSON(http.StatusOK, result)

}

// UpdateConsumable godoc
//
//	@Security		BearerAuth
//	@Summary		Atualizar consumível
//	@Description	Atualiza os dados de um consumível existente
//	@Tags			Consumíveis
//	@Accept			json
//	@Produce		json
//	@Param			uid		path		string					true	"UID do consumível"
//	@Param			request	body		UpdateConsumableRequest	true	"Dados para atualização"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/consumables/{uid} [put]
func (h *Handler) UpdateConsumable(c *gin.Context) {
	var (
		parentCtx = tracing.NewRelicContext(c)
		in        = UpdateConsumableRequest{}
		uid       = c.Param("uid")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-update-consumable")
	defer func() {
		span.End()
	}()

	if uid == "" {
		c.JSON(http.StatusBadRequest, "uid é obrigatorio")
		return
	}

	_, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, "uid invalido")
		return
	}

	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	span.AddEvent("http-handler-update-consumable-params", tracing.Attributes{
		"model": in,
	})

	err = h.services.UpdateConsumable(ctx, cports.UpdateConsumable{
		UID:         null.StringFrom(uid),
		Name:        null.StringFrom(in.Name),
		Price:       null.FloatFrom(in.Price),
		Description: null.StringFrom(in.Description),
	})
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-update-consumable-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})

		return
	}

	span.AddEvent("http-handler-update-consumable-success", nil)

	c.JSON(http.StatusNoContent, nil)

}

// DeleteConsumable godoc
//
//	@Security		BearerAuth
//	@Summary		Deletar consumível
//	@Description	Remove um consumível pelo seu identificador único
//	@Tags			Consumíveis
//	@Produce		json
//	@Param			uid	path		string	true	"UID do consumível"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/consumables/{uid} [delete]
func (h *Handler) DeleteConsumable(c *gin.Context) {

	var (
		parentCtx = tracing.NewRelicContext(c)
		uid       = c.Param("uid")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-delete-consumable")
	defer func() {
		span.End()
	}()

	if uid == "" {
		c.JSON(http.StatusBadRequest, "uid é obrigatorio")
		return
	}

	_, err := uuid.Parse(uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, "uid invalido")
		return
	}

	span.AddEvent("http-handler-delete-consumable-params", tracing.Attributes{
		"model": uid,
	})

	err = h.services.DeleteConsumable(ctx, uid)
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-delete-consumable-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})

		return
	}

	span.AddEvent("http-handler-delete-consumable-success", nil)

	c.JSON(http.StatusNoContent, nil)

}
