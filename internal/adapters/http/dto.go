package http

import "time"

// --- Consumable DTOs ---

// SaveConsumableRequest representa o body para criação de consumível
type SaveConsumableRequest struct {
	Name        string  `json:"name" binding:"required" example:"Café Especial"`
	Description string  `json:"description" example:"Café arábica torrado"`
	Price       float64 `json:"price" binding:"required" example:"29.90"`
	Code        string  `json:"code" binding:"required" example:"CAF-001"`
	Quantity    int32   `json:"quantity" binding:"required" example:"100"`
}

// UpdateConsumableRequest representa o body para atualização de consumível
type UpdateConsumableRequest struct {
	Name        string  `json:"name" binding:"required" example:"Café Especial Premium"`
	Description string  `json:"description" example:"Café arábica torrado premium"`
	Price       float64 `json:"price" binding:"required" example:"39.90"`
}

// --- Service Definition DTOs ---

// SaveServiceDefinitionRequest representa o body para criação de definição de serviço
type SaveServiceDefinitionRequest struct {
	Name        string  `json:"name" binding:"required" example:"Limpeza Geral"`
	Description string  `json:"description" example:"Serviço de limpeza completa"`
	Price       float64 `json:"price" binding:"required" example:"150.00"`
}

// UpdateServiceDefinitionRequest representa o body para atualização de definição de serviço
type UpdateServiceDefinitionRequest struct {
	Name        string  `json:"name" binding:"required" example:"Limpeza Geral Premium"`
	Description string  `json:"description" example:"Serviço de limpeza completa premium"`
	Price       float64 `json:"price" binding:"required" example:"200.00"`
}

// --- Response DTOs ---

// MessageResponse resposta genérica com mensagem
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse resposta de erro
type ErrorResponse struct {
	Code  int    `json:"code" example:"400"`
	Error string `json:"error" example:"invalid request"`
}

// PaginatedResponse resposta paginada genérica
type PaginatedResponse struct {
	Data  any   `json:"data"`
	Total int64 `json:"total" example:"100"`
	Page  int64 `json:"page" example:"1"`
	Limit int64 `json:"limit" example:"10"`
}

type GetServiceDefinitionResponse struct {
	UID         string    `json:"uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"Limpeza Geral"`
	Description string    `json:"description" example:"Serviço de limpeza completa"`
	Price       float64   `json:"price" example:"150.00"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}

type GetConsumableResponse struct {
	UID         string    `json:"uid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"Café Especial"`
	Description string    `json:"description" example:"Café arábica torrado"`
	Code        string    `json:"code" example:"CAF-001"`
	Quantity    int64     `json:"quantity" example:"100"`
	Price       float64   `json:"price" example:"29.90"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T12:00:00Z"`
}
