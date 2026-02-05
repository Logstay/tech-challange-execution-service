package http

import (
	services "github.com/Logstay/tech-challange-wherehouse-service/internal/interfaces/services"
)

type Handler struct {
	services services.ServiceFactory
}

func NewHandler(services services.ServiceFactory) *Handler {
	return &Handler{
		services: services,
	}
}
