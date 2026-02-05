package http

import bPorts "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/billing"

type CardPaymentRequest struct {
	Installments      int              `json:"installments" example:"1"`
	TransactionAmount float64          `json:"transaction_amount" example:"100.00"`
	CardToken         string           `json:"token" example:"ea1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"`
	Description       string           `json:"description" example:"Pagamento de serviço automotivo"`
	PaymentMethodID   string           `json:"payment_method_id" example:"visa"`
	Payer             CardPayerRequest `json:"payer"`
	BodyCard          BodyCard         `json:"card_details"`
	ExternalReference string           `json:"external_reference,omitempty" example:"ORDER-12345"`
	NotificationURL   string           `json:"notification_url,omitempty" example:"https://seu-dominio.com/webhook/mercadopago"`
}

type CardPayerRequest struct {
	Email          string                 `json:"email" example:"customer@example.com"`
	Identification *IdentificationRequest `json:"identification"`
}

type IdentificationRequest struct {
	Type   string `json:"type" example:"CPF"`
	Number string `json:"number" example:"12345678909"`
}

type CardPaymentResponse struct {
	ID                int64     `json:"id" example:"123456789"`
	Status            string    `json:"status" example:"approved"`
	StatusDetail      string    `json:"status_detail" example:"accredited"`
	TransactionAmount float64   `json:"transaction_amount" example:"100.00"`
	Description       string    `json:"description" example:"Pagamento de serviço automotivo"`
	ExternalReference string    `json:"external_reference,omitempty" example:"ORDER-12345"`
	PaymentMethodID   string    `json:"payment_method_id" example:"visa"`
	PaymentTypeID     string    `json:"payment_type_id" example:"credit_card"`
	Installments      int       `json:"installments" example:"1"`
	Card              *CardInfo `json:"card,omitempty"`
}

type CardInfo struct {
	FirstSixDigits string `json:"first_six_digits" example:"450995"`
	LastFourDigits string `json:"last_four_digits" example:"3704"`
}

type BodyCard struct {
	SiteID          string     `json:"site_id" example:"MLB"`
	CardNumber      string     `json:"card_number" example:"4235647728025682"`
	ExpirationYear  string     `json:"expiration_year" example:"2030"`
	ExpirationMonth string     `json:"expiration_month" example:"11"`
	SecurityCode    string     `json:"security_code" example:"123"`
	Cardholder      Cardholder `json:"cardholder"`
}

type Cardholder struct {
	Identification Identification `json:"identification"`
	Name           string         `json:"name" example:"APRO"`
}

type Identification struct {
	Type   string `json:"type" example:"CPF"`
	Number string `json:"number" example:"12345678909"`
}

type QuotationResponse struct {
	ID                string                         `json:"id"`
	PaymentLink       string                         `json:"payment_link"`
	SandboxLink       string                         `json:"sandbox_link"`
	ExternalReference string                         `json:"external_reference"`
	TotalAmount       float64                        `json:"total_amount"`
	Currency          string                         `json:"currency"`
	Status            string                         `json:"status"`
	CreatedAt         string                         `json:"created_at"`
	ExpiresAt         string                         `json:"expires_at"`
	Items             []bPorts.PreferenceItemRequest `json:"items"`
	Metadata          map[string]interface{}         `json:"metadata,omitempty"`
}
