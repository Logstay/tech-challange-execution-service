package http

import "time"

type SaveConsumableRequest struct {
	Name        string  `json:"name" example:"Óleo de Motor GTX"`
	Price       float64 `json:"price" example:"20"`
	Quantity    int64   `json:"quantity" example:"10"`
	Description string  `json:"description" example:"produto de alta densidade"`
}

type GetConsumableResponse struct {
	UID         string     `json:"uid"  example:"8c2a4b4a-6c9d-4e8e-8c2a-4b4a6c9d4e8e"`
	Name        string     `json:"name,omitempty"  example:"Óleo de Motor GTX"`
	Price       string     `json:"price,omitempty"  example:"20"`
	Code        string     `json:"code,omitempty"  example:"cECAZFNnmq5uM2gpST"`
	Quantity    int64      `json:"quantity,omitempty"  example:"10"`
	Description string     `json:"description,omitempty"  example:"produto de alta densidade"`
	CreatedAt   *time.Time `json:"created_at,omitempty"  example:"2022-04-29T13:39:50.61744Z"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"  example:"2022-04-29T13:39:50.61744Z"`
}

type UpdateConsumableRequest struct {
	Name        string  `json:"name" example:"Óleo Mobil"`
	Price       float64 `json:"price" example:"22"`
	Description string  `json:"description" example:"produto de baixa densidade"`
}

type ListConsumableResponse struct {
	Data       []GetConsumableResponse `json:"data"`
	Pagination Pagination              `json:"pagination"`
}

type SaveServiceDefinitionRequest struct {
	Name        string  `json:"name" example:"Alinhamento e Balanceamento"`
	Price       float64 `json:"price" example:"200"`
	Description string  `json:"description" example:"Serviço completo de alinhamento e balanceamento de rodas com equipamento de alta precisão"`
}

type GetServiceDefinitionResponse struct {
	UID         string     `json:"uid"  example:"8c2a4b4a-6c9d-4e8e-8c2a-4b4a6c9d4e8e"`
	Name        string     `json:"name,omitempty"  example:"Alinhamento e Balanceamento"`
	Price       string     `json:"price,omitempty"  example:"200"`
	Description string     `json:"description,omitempty" example:"Serviço completo de alinhamento e balanceamento de rodas com equipamento de alta precisão"`
	CreatedAt   *time.Time `json:"created_at,omitempty"  example:"2022-04-29T13:39:50.61744Z"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"  example:"2022-04-29T13:39:50.61744Z"`
}

type UpdateServiceDefinitionRequest struct {
	Name        string  `json:"name" example:"Troca de Óleo"`
	Price       float64 `json:"price" example:"75"`
	Description string  `json:"description" example:"Serviço completo de alinhamento e balanceamento de rodas com equipamento de alta precisão"`
}

type ListServiceDefinitionResponse struct {
	Data       []GetServiceDefinitionResponse `json:"data"`
	Pagination Pagination                     `json:"pagination"`
}

type Pagination struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
}
