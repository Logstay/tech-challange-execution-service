package http

import (
	"net/http"

	sdports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SaveServiceDefinition godoc
//
//	@Security		BearerAuth
//	@Summary		Criar definição de serviço
//	@Description	Cria uma nova definição de serviço
//	@Tags			Definições de Serviço
//	@Accept			json
//	@Produce		json
//	@Param			request	body		SaveServiceDefinitionRequest	true	"Dados da definição de serviço"
//	@Success		201		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/service-definitions [post]
func (h *Handler) SaveServiceDefinition(c *gin.Context) {

	var (
		parentCtx = tracing.NewRelicContext(c)
		in        = SaveServiceDefinitionRequest{}
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-save-serviceDefinition")
	defer func() {
		span.End()
	}()

	if err := c.ShouldBind(&in); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	span.AddEvent("http-handler-save-serviceDefinition-params", tracing.Attributes{
		"model": in,
	})

	err := h.services.SaveServiceDefinition(ctx, sdports.SaveServiceDefinition{
		Name:  null.StringFrom(in.Name),
		Price: null.FloatFrom(in.Price),
		Description: null.StringFrom(
			in.Description,
		),
	})
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-save-serviceDefinition-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})

		return
	}

	span.AddEvent("http-handler-save-serviceDefinition-success", nil)

	c.JSON(http.StatusCreated, nil)

}

// GetServiceDefinition godoc
//
//	@Security		BearerAuth
//	@Summary		Buscar definição de serviço por UID
//	@Description	Retorna uma definição de serviço pelo seu identificador único
//	@Tags			Definições de Serviço
//	@Produce		json
//	@Param			uid	path		string	true	"UID da definição de serviço"
//	@Success		200	{object}	GetServiceDefinitionResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/service-definitions/{uid} [get]
func (h *Handler) GetServiceDefinition(c *gin.Context) {

	var (
		parentCtx = tracing.NewRelicContext(c)
		uid       = c.Param("uid")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-get-serviceDefinition")
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

	resp, err := h.services.GetServiceDefinitionByUID(ctx, uid)
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-get-serviceDefinition-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})
		return
	}

	result := GetServiceDefinitionResponse{
		UID:         resp.UID.String,
		Name:        resp.Name.String,
		Price:       float64ToReais(resp.Price.Float64),
		Description: resp.Description.String,
		CreatedAt:   resp.CreatedAt.Ptr(),
		UpdatedAt:   resp.UpdatedAt.Ptr(),
	}

	span.AddEvent("http-handler-get-serviceDefinition-success", nil)

	c.JSON(http.StatusOK, result)

}

// ListServiceDefinition godoc
//
//	@Security		BearerAuth
//	@Summary		Listar definições de serviço
//	@Description	Retorna uma lista paginada de definições de serviço
//	@Tags			Definições de Serviço
//	@Produce		json
//	@Param			page	query		int		false	"Número da página"	default(0)
//	@Param			limit	query		int		false	"Itens por página"	default(10)
//	@Param			name	query		string	false	"Filtrar por nome"
//	@Success		200		{object}	PaginatedResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/service-definitions [get]
func (h *Handler) ListServiceDefinition(c *gin.Context) {

	page, limit := pageFromRequest(c.Request)
	var (
		parentCtx = tracing.NewRelicContext(c)
		name      = c.Query("name")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-list-serviceDefinition")
	defer func() {
		span.End()
	}()

	resp, total, err := h.services.ListServiceDefinition(ctx, sdports.Filter{
		Name:  null.StringFrom(name),
		Page:  null.IntFrom(page),
		Limit: null.IntFrom(limit),
	})
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-list-serviceDefinition-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})
		return
	}

	result := ListServiceDefinitionResponse{}
	for _, r := range resp {
		result.Data = append(result.Data, GetServiceDefinitionResponse{
			UID:         r.UID.String,
			Name:        r.Name.String,
			Price:       float64ToReais(r.Price.Float64),
			Description: r.Description.String,
			CreatedAt:   r.CreatedAt.Ptr(),
			UpdatedAt:   r.UpdatedAt.Ptr(),
		})
	}

	result.Pagination = Pagination{
		Total: total,
		Limit: limit,
		Page:  page,
	}

	span.AddEvent("http-handler-list-serviceDefinition-success", nil)

	c.JSON(http.StatusOK, result)

}

// UpdateServiceDefinition godoc
//
//	@Security		BearerAuth
//	@Summary		Atualizar definição de serviço
//	@Description	Atualiza os dados de uma definição de serviço existente
//	@Tags			Definições de Serviço
//	@Accept			json
//	@Produce		json
//	@Param			uid		path		string							true	"UID da definição de serviço"
//	@Param			request	body		UpdateServiceDefinitionRequest	true	"Dados para atualização"
//	@Success		200		{object}	MessageResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		401		{object}	ErrorResponse
//	@Failure		404		{object}	ErrorResponse
//	@Failure		500		{object}	ErrorResponse
//	@Router			/service-definitions/{uid} [put]
func (h *Handler) UpdateServiceDefinition(c *gin.Context) {

	var (
		parentCtx = tracing.NewRelicContext(c)
		in        = UpdateServiceDefinitionRequest{}
		uid       = c.Param("uid")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-update-serviceDefinition")
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

	span.AddEvent("http-handler-update-serviceDefinition-params", tracing.Attributes{
		"model": in,
	})

	err = h.services.UpdateServiceDefinition(ctx, sdports.UpdateServiceDefinition{
		UID:         null.StringFrom(uid),
		Name:        null.StringFrom(in.Name),
		Price:       null.FloatFrom(in.Price),
		Description: null.StringFrom(in.Description),
	})
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-update-serviceDefinition-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})

		return
	}

	span.AddEvent("http-handler-update-serviceDefinition-success", nil)

	c.JSON(http.StatusNoContent, nil)

}

// DeleteServiceDefinition godoc
//
//	@Security		BearerAuth
//	@Summary		Deletar definição de serviço
//	@Description	Remove uma definição de serviço pelo seu identificador único
//	@Tags			Definições de Serviço
//	@Produce		json
//	@Param			uid	path		string	true	"UID da definição de serviço"
//	@Success		200	{object}	MessageResponse
//	@Failure		401	{object}	ErrorResponse
//	@Failure		404	{object}	ErrorResponse
//	@Failure		500	{object}	ErrorResponse
//	@Router			/service-definitions/{uid} [delete]
func (h *Handler) DeleteServiceDefinition(c *gin.Context) {

	var (
		parentCtx = tracing.NewRelicContext(c)
		uid       = c.Param("uid")
	)

	ctx, span := h.tracer.NewSpan(parentCtx, "http-handler-delete-serviceDefinition")
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

	span.AddEvent("http-handler-delete-serviceDefinition-params", tracing.Attributes{
		"model": uid,
	})

	err = h.services.DeleteServiceDefinition(ctx, uid)
	if err != nil {
		handlerResp := httperr.ErrorHandler(err, c)
		span.AddEvent("http-handler-delete-serviceDefinition-error", tracing.Attributes{
			"err":             handlerResp.Error,
			"http.statusCode": handlerResp.Code,
		})

		return
	}

	span.AddEvent("http-handler-delete-serviceDefinition-success", nil)

	c.JSON(http.StatusOK, nil)

}
