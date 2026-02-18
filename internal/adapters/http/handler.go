package http

import (
	"net/http"
	"strconv"

	services "github.com/Logstay/tech-challange-execution-service/internal/interfaces/services"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Handler struct {
	services services.ServiceFactory
	tracer   *tracing.Tracer
}

func NewHandler(services services.ServiceFactory, tracer *tracing.Provider) *Handler {
	return &Handler{
		services: services,
		tracer:   tracer.NewTracer("handlers"),
	}
}

func pageFromRequest(r *http.Request) (page int64, limit int64) {
	limit = 10

	if p := r.URL.Query().Get("page"); p != "" {
		if val, err := strconv.ParseInt(p, 10, 64); err == nil {
			page = val
		}
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		if val, err := strconv.ParseInt(l, 10, 64); err == nil {
			limit = val
		}
	}

	return
}

func float64ToReais(valor float64) (value string) {
	if valor == 0 {
		return
	}

	p := message.NewPrinter(language.BrazilianPortuguese)

	value = p.Sprintf("R$ %.2f", valor)

	return
}
