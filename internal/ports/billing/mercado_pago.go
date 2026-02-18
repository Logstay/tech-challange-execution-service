package billing

import (
	"context"

	"github.com/mercadopago/sdk-go/pkg/preference"
)

type MercadoPagoBilling interface {
	GenerateCardToken(ctx context.Context, request BodyCard) (cardTokenID string, err error)
	CreateCardPayment(ctx context.Context, request CardPaymentRequest) (resp *CardPaymentResponse, err error)
	GetPayment(ctx context.Context, paymentID int64) (resp *CardPaymentResponse, err error)
	CreateRefund(ctx context.Context, paymentID int64, amount *float64) (resp *RefundResponse, err error)
	CancelPayment(ctx context.Context, paymentID int64) (err error)
	SearchPayments(ctx context.Context, filters PaymentSearchFilters) (resp *PaymentSearchResponse, err error)
	CreatePreference(ctx context.Context, request PreferenceRequest) (resp *PreferenceResponse, err error)
	GetPreference(ctx context.Context, preferenceID string) (resp *preference.Response, err error)
	UpdatePreference(ctx context.Context, preferenceID string, request PreferenceRequest) (resp *PreferenceResponse, err error)
}

type OrderRequest struct {
	Type              string              `json:"type"`
	TotalAmount       string              `json:"total_amount"`
	ExternalReference string              `json:"external_reference"`
	Currency          string              `json:"currency"`
	Transactions      *TransactionRequest `json:"transactions"`
	Payer             *PayerRequest       `json:"payer"`
}

type TravelPassengerRequest struct {
	FirstName      string                 `json:"first_name,omitempty"`
	LastName       string                 `json:"last_name,omitempty"`
	Identification *IdentificationRequest `json:"identification,omitempty"`
}

type TravelRouteRequest struct {
	Departure         string `json:"departure,omitempty"`
	Destination       string `json:"destination,omitempty"`
	DepartureDateTime string `json:"departure_date_time,omitempty"`
	ArrivalDateTime   string `json:"arrival_date_time,omitempty"`
	Company           string `json:"company,omitempty"`
}

type AdditionalInfoRequest struct {
	PayerAuthenticationType            string                    `json:"payer.authentication_type,omitempty"`
	PayerRegistrationDate              string                    `json:"payer.registration_date,omitempty"`
	PayerIsPrimeUser                   bool                      `json:"payer.is_prime_user,omitempty"`
	PayerIsFirstPurchaseOnLine         bool                      `json:"payer.is_first_purchase_online,omitempty"`
	PayerLastPurchase                  string                    `json:"payer.last_purchase,omitempty"`
	ShipmentExpress                    bool                      `json:"shipment.express,omitempty"`
	ShipmentLocalPickup                bool                      `json:"shipment.local_pickup,omitempty"`
	PlatFormShipmentDeliveryPromise    string                    `json:"platform.shipment.delivery_promise,omitempty"`
	PlatFormShipmentDropShipping       string                    `json:"platform.shipment.drop_shipping,omitempty"`
	PlatformShipmentSafety             string                    `json:"platform.shipment.safety,omitempty"`
	PlatformShipmentTrackingCode       string                    `json:"platform.shipment.tracking.code,omitempty"`
	PlatformShipmentTrackingStatus     string                    `json:"platform.shipment.tracking.status,omitempty"`
	PlatformShipmentWithdrawn          bool                      `json:"platform.shipment.withdrawn,omitempty"`
	PlatformSellerID                   string                    `json:"platform.seller.id,omitempty"`
	PlatformSellerName                 string                    `json:"platform.seller.name,omitempty"`
	PlatformSellerEmail                string                    `json:"platform.seller.email,omitempty"`
	PlatformSellerStatus               string                    `json:"platform.seller.status,omitempty"`
	PlatformSellerReferralURL          string                    `json:"platform.seller.referral_url,omitempty"`
	PlatformSellerRegistrationDate     string                    `json:"platform.seller.registration_date,omitempty"`
	PlatformSellerHiredPlan            string                    `json:"platform.seller.hired_plan,omitempty"`
	PlatformSellerBusinessType         string                    `json:"platform.seller.business_type,omitempty"`
	PlatformSellerAddressZipCode       string                    `json:"platform.seller.address.zip_code,omitempty"`
	PlatformSellerAddressStreetName    string                    `json:"platform.seller.address.street_name,omitempty"`
	PlatformSellerAddressStreetNumber  string                    `json:"platform.seller.address.street_number,omitempty"`
	PlatformSellerAddressCity          string                    `json:"platform.seller.address.city,omitempty"`
	PlatformSellerAddressState         string                    `json:"platform.seller.address.state,omitempty"`
	PlatformSellerAddressComplement    string                    `json:"platform.seller.address.complement,omitempty"`
	PlatformSellerAddressCountry       string                    `json:"platform.seller.address.country,omitempty"`
	PlatformSellerIdentificationType   string                    `json:"platform.seller.identification.type,omitempty"`
	PlatformSellerIdentificationNumber string                    `json:"platform.seller.identification.number,omitempty"`
	PlatformSellerPhoneAreaCode        string                    `json:"platform.seller.phone.area_code,omitempty"`
	PlatformSellerPhoneNumber          string                    `json:"platform.seller.phone.number,omitempty"`
	PlatformAuthentication             string                    `json:"platform.authentication,omitempty"`
	TravelPassengers                   *[]TravelPassengerRequest `json:"travel.passengers,omitempty"`
	TravelRoutes                       *[]TravelRouterRequest    `json:"travel.routes,omitempty"`
}

type TravelRouterRequest struct {
	Departure         string `json:"departure,omitempty"`
	Destination       string `json:"destination,omitempty"`
	DepartureDateTime string `json:"departure_date_time,omitempty"`
	ArrivalDateTime   string `json:"arrival_date_time,omitempty"`
	Company           string `json:"company,omitempty"`
}
type PayerAddressRequest struct {
	ZipCode      string `json:"zip_code,omitempty"`
	StreetName   string `json:"street_name,omitempty"`
	StreetNumber string `json:"street_number,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Complement   string `json:"complement,omitempty"`
	Country      string `json:"country,omitempty"`
}

type ShipmentRequest struct {
	Address *AddressRequest `json:"address,omitempty"`
}

type TransactionRequest struct {
	Payments []PaymentRequest `json:"payments,omitempty"`
}

type PaymentRequest struct {
	Amount            string                    `json:"amount,omitempty"`
	ExpirationTime    string                    `json:"expiration_time,omitempty"`
	PaymentMethod     *PaymentMethodRequest     `json:"payment_method,omitempty"`
	AutomaticPayments *AutomaticPaymentsRequest `json:"automatic_payments,omitempty"`
	StoredCredential  *StoredCredentialRequest  `json:"stored_credential,omitempty"`
	SubscriptionData  *SubscriptionDataRequest  `json:"subscription_data,omitempty"`
}

type PaymentMethodRequest struct {
	ID                  string `json:"id,omitempty"`
	Type                string `json:"type,omitempty"`
	Token               string `json:"token,omitempty"`
	StatementDescriptor string `json:"statement_descriptor,omitempty"`
	Installments        int    `json:"installments,omitempty"`
}

type AutomaticPaymentsRequest struct {
	PaymentProfileID string `json:"payment_profile_id,omitempty"`
	ScheduleDate     string `json:"schedule_date,omitempty"`
	DueDate          string `json:"due_date,omitempty"`
	Retries          int    `json:"retries,omitempty"`
}

type StoredCredentialRequest struct {
	PaymentInitiator   string `json:"payment_initiator,omitempty"`
	Reason             string `json:"reason,omitempty"`
	StorePaymentMethod bool   `json:"store_payment_method,omitempty"`
	FirstPayment       bool   `json:"first_payment,omitempty"`
}

type SubscriptionDataRequest struct {
	InvoiceID            string                       `json:"invoice_id,omitempty"`
	BillingDate          string                       `json:"billing_date,omitempty"`
	SubscriptionSequence *SubscriptionSequenceRequest `json:"subscription_sequence,omitempty"`
	InvoicePeriod        *InvoicePeriodRequest        `json:"invoice_period,omitempty"`
}

type SubscriptionSequenceRequest struct {
	Number int `json:"number,omitempty"`
	Total  int `json:"total,omitempty"`
}

type InvoicePeriodRequest struct {
	Type   string `json:"type,omitempty"`
	Period int    `json:"period,omitempty"`
}

type PayerRequest struct {
	Email          string                 `json:"email,omitempty"`
	FirstName      string                 `json:"first_name,omitempty"`
	LastName       string                 `json:"last_name,omitempty"`
	CustomerID     string                 `json:"customer_id,omitempty"`
	EntityType     string                 `json:"entity_type,omitempty"`
	Identification *IdentificationRequest `json:"identification,omitempty"`
	Phone          *PhoneRequest          `json:"phone,omitempty"`
	Address        *PayerAddressRequest   `json:"address,omitempty"`
}

type PayerAddress struct {
	ZipCode      string `json:"zip_code,omitempty"`
	StreetName   string `json:"street_name,omitempty"`
	StreetNumber string `json:"street_number,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
}

type IdentificationRequest struct {
	Type   string `json:"type,omitempty"`
	Number string `json:"number,omitempty"`
}

type PhoneRequest struct {
	AreaCode string `json:"area_code,omitempty"`
	Number   string `json:"number,omitempty"`
}

type AddressRequest struct {
	StreetName   string `json:"street_name,omitempty"`
	StreetNumber string `json:"street_number,omitempty"`
	ZipCode      string `json:"zip_code,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	State        string `json:"state,omitempty"`
	City         string `json:"city,omitempty"`
	Complement   string `json:"complement,omitempty"`
}

type ItemsRequest struct {
	Title        string `json:"title,omitempty"`
	Type         string `json:"type,omitempty"`
	Warranty     bool   `json:"warranty,omitempty"`
	EventDate    string `json:"event_date,omitempty"`
	UnitPrice    string `json:"unit_price,omitempty"`
	ExternalCode string `json:"external_code,omitempty"`
	CategoryID   string `json:"category_id,omitempty"`
	Description  string `json:"description,omitempty"`
	PictureURL   string `json:"picture_url,omitempty"`
	Quantity     int    `json:"quantity,omitempty"`
}

type RefundRequest struct {
	Transactions []RefundTransaction `json:"transactions,omitempty"`
}

type RefundTransaction struct {
	ID     string `json:"id,omitempty"`
	Amount string `json:"amount,omitempty"`
}

type ConfigRequest struct {
	PaymentMethod *PaymentMethodConfigRequest `json:"payment_method,omitempty"`
	Online        *OnlineConfigRequest        `json:"online,omitempty"`
}

type PaymentMethodConfigRequest struct {
	NotAllowedIDs       []string `json:"not_allowed_ids,omitempty"`
	NotAllowedTypes     []string `json:"not_allowed_types,omitempty"`
	DefaultID           string   `json:"default_id,omitempty"`
	MaxInstallments     int      `json:"max_installments,omitempty"`
	DefaultInstallments int      `json:"default_installments,omitempty"`
}

type OnlineConfigRequest struct {
	CallbackURL         string                      `json:"callback_url,omitempty"`
	SuccessURL          string                      `json:"success_url,omitempty"`
	PendingURL          string                      `json:"pending_url,omitempty"`
	FailureURL          string                      `json:"failure_url,omitempty"`
	AutoReturnURL       string                      `json:"auto_return_url,omitempty"`
	DifferentialPricing *DifferentialPricingRequest `json:"differential_pricing,omitempty"`
	TransactionSecurity *TransactionSecurityRequest `json:"transaction_security,omitempty"`
}

type DifferentialPricingRequest struct {
	ID int `json:"id,omitempty"`
}

type TransactionSecurityRequest struct {
	Validation     string `json:"validation,omitempty"`
	LiabilityShift string `json:"liability_shift,omitempty"`
}

type PixPaymentRequest struct {
	TransactionAmount float64                   `json:"transaction_amount"`
	Description       string                    `json:"description"`
	PaymentMethodID   string                    `json:"payment_method_id"`
	Payer             PixPayerRequest           `json:"payer"`
	ExternalReference string                    `json:"external_reference,omitempty"`
	NotificationURL   string                    `json:"notification_url,omitempty"`
	AdditionalInfo    *PixAdditionalInfoRequest `json:"additional_info,omitempty"`
}

type PixPayerRequest struct {
	Email          string                  `json:"email"`
	FirstName      string                  `json:"first_name,omitempty"`
	LastName       string                  `json:"last_name,omitempty"`
	Identification *IdentificationRequest  `json:"identification,omitempty"`
	Address        *PixPayerAddressRequest `json:"address,omitempty"`
}

type PixPayerAddressRequest struct {
	ZipCode      string `json:"zip_code,omitempty"`
	StreetName   string `json:"street_name,omitempty"`
	StreetNumber string `json:"street_number,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
}

type PixAdditionalInfoRequest struct {
	Items []PixItemRequest `json:"items,omitempty"`
}

type PixItemRequest struct {
	ID          string  `json:"id,omitempty"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	PictureURL  string  `json:"picture_url,omitempty"`
	CategoryID  string  `json:"category_id,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
	UnitPrice   float64 `json:"unit_price,omitempty"`
}

type PixPaymentResponse struct {
	ID                 int64                  `json:"id"`
	Status             string                 `json:"status"`
	StatusDetail       string                 `json:"status_detail"`
	TransactionAmount  float64                `json:"transaction_amount"`
	Description        string                 `json:"description"`
	ExternalReference  string                 `json:"external_reference,omitempty"`
	QRCode             string                 `json:"qr_code"`
	QRCodeBase64       string                 `json:"qr_code_base64"`
	ExpirationDate     string                 `json:"expiration_date,omitempty"`
	PointOfInteraction *PixPointOfInteraction `json:"point_of_interaction,omitempty"`
}

type PixPointOfInteraction struct {
	Type            string              `json:"type"`
	TransactionData *PixTransactionData `json:"transaction_data,omitempty"`
}

type PixTransactionData struct {
	QRCode       string `json:"qr_code"`
	QRCodeBase64 string `json:"qr_code_base64"`
	TicketURL    string `json:"ticket_url,omitempty"`
}

type PreferenceRequest struct {
	Items             []PreferenceItemRequest `json:"items"`
	ExternalReference string                  `json:"external_reference,omitempty"`
	NotificationURL   string                  `json:"notification_url,omitempty"`
	BackURLs          *BackURLsRequest        `json:"back_urls,omitempty"`
	AutoReturn        string                  `json:"auto_return,omitempty"`
	Payer             *PreferencePayerRequest `json:"payer,omitempty"`
	// Campos para Orçamento/Quotation com vigência
	Expires             bool                   `json:"expires,omitempty"`
	ExpirationDateFrom  *string                `json:"expiration_date_from,omitempty"` // ISO 8601: 2024-02-04T10:00:00.000-03:00
	ExpirationDateTo    *string                `json:"expiration_date_to,omitempty"`   // ISO 8601: 2024-02-11T23:59:59.000-03:00
	DateOfExpiration    *string                `json:"date_of_expiration,omitempty"`   // Para pagamentos offline (PIX, boleto)
	StatementDescriptor string                 `json:"statement_descriptor,omitempty"`
	Metadata            map[string]interface{} `json:"metadata,omitempty"`
}

type PreferenceItemRequest struct {
	Title       string  `json:"title"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	CurrencyID  string  `json:"currency_id,omitempty"`
	Description string  `json:"description,omitempty"`
}

type BackURLsRequest struct {
	Success string `json:"success,omitempty"`
	Failure string `json:"failure,omitempty"`
	Pending string `json:"pending,omitempty"`
}

type PreferencePayerRequest struct {
	Name    string        `json:"name,omitempty"`
	Surname string        `json:"surname,omitempty"`
	Email   string        `json:"email,omitempty"`
	Phone   *PhoneRequest `json:"phone,omitempty"`
}

type PreferenceResponse struct {
	ID                 string                  `json:"id"`
	InitPoint          string                  `json:"init_point"`         // URL para pagamento em produção
	SandboxInitPoint   string                  `json:"sandbox_init_point"` // URL para testes
	Items              []PreferenceItemRequest `json:"items"`
	ExternalReference  string                  `json:"external_reference,omitempty"`
	DateCreated        string                  `json:"date_created,omitempty"`
	ExpirationDateFrom *string                 `json:"expiration_date_from,omitempty"`
	ExpirationDateTo   *string                 `json:"expiration_date_to,omitempty"`
	DateOfExpiration   *string                 `json:"date_of_expiration,omitempty"`
	Expires            bool                    `json:"expires,omitempty"`
	CollectorID        int64                   `json:"collector_id,omitempty"`
	ClientID           string                  `json:"client_id,omitempty"`
	Marketplace        string                  `json:"marketplace,omitempty"`
	Status             string                  `json:"status,omitempty"`
	Metadata           map[string]interface{}  `json:"metadata,omitempty"`
}

type CardPaymentRequest struct {
	Installments      int              `json:"installments"`
	TransactionAmount float64          `json:"transaction_amount"`
	CardToken         string           `json:"token"`
	Description       string           `json:"description"`
	PaymentMethodID   string           `json:"payment_method_id"`
	Payer             CardPayerRequest `json:"payer"`
	BodyCard          BodyCard         `json:"card_details"`
	ExternalReference string           `json:"external_reference,omitempty"`
	NotificationURL   string           `json:"notification_url,omitempty"`
}

type CardPayerRequest struct {
	Email          string                 `json:"email"`
	Identification *IdentificationRequest `json:"identification"`
}

type CardPaymentResponse struct {
	ID                int64     `json:"id"`
	Status            string    `json:"status"`
	StatusDetail      string    `json:"status_detail"`
	TransactionAmount float64   `json:"transaction_amount"`
	Description       string    `json:"description"`
	ExternalReference string    `json:"external_reference,omitempty"`
	PaymentMethodID   string    `json:"payment_method_id"`
	PaymentTypeID     string    `json:"payment_type_id"`
	Installments      int       `json:"installments"`
	Card              *CardInfo `json:"card,omitempty"`
}

type CardInfo struct {
	FirstSixDigits string `json:"first_six_digits"`
	LastFourDigits string `json:"last_four_digits"`
}

type BodyCard struct {
	SiteID          string     `json:"site_id"`
	CardNumber      string     `json:"card_number"`
	ExpirationYear  string     `json:"expiration_year"`
	ExpirationMonth string     `json:"expiration_month"`
	SecurityCode    string     `json:"security_code"`
	Cardholder      Cardholder `json:"cardholder"`
}

type Cardholder struct {
	Identification Identification `json:"identification"`
	Name           string         `json:"name"`
}

type Identification struct {
	Type   string `json:"type"`
	Number string `json:"number"`
}

type RefundResponse struct {
	ID                   int64   `json:"id"`
	PaymentID            int64   `json:"payment_id"`
	Amount               float64 `json:"amount"`
	Status               string  `json:"status"`
	DateCreated          string  `json:"date_created"`
	SourceID             int64   `json:"source_id"`
	RefundMode           string  `json:"refund_mode"`
	AdjustmentAmount     float64 `json:"adjustment_amount"`
	StatusDetail         string  `json:"status_detail,omitempty"`
	UniqueSequenceNumber string  `json:"unique_sequence_number,omitempty"`
}

type PaymentSearchFilters struct {
	ExternalReference string `json:"external_reference,omitempty"`
	Status            string `json:"status,omitempty"`
	DateFrom          string `json:"date_from,omitempty"`
	DateTo            string `json:"date_to,omitempty"`
	Limit             int    `json:"limit,omitempty"`
	Offset            int    `json:"offset,omitempty"`
}

type PaymentSearchResponse struct {
	Paging  *PagingInfo           `json:"paging"`
	Results []CardPaymentResponse `json:"results"`
}

type PagingInfo struct {
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
