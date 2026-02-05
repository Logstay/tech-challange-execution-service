package mercadopago

import (
	"context"
	"errors"
	"fmt"

	bPorts "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/billing"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/logger"
	"github.com/mercadopago/sdk-go/pkg/cardtoken"
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
	"github.com/mercadopago/sdk-go/pkg/refund"
)

type mercadoPagoPayment struct {
	config *config.Config
	log    logger.Logger
}

func NewMercadoPagoPayment(accessToken string, logger logger.Logger) (payment bPorts.MercadoPagoBilling, err error) {

	if accessToken == "" {
		err = errors.New("access token is required")
		logger.Error("Mercado Pago access token is missing", err, nil)
		return
	}

	c, err := config.New(accessToken)
	if err != nil {
		err = errors.New("failed to create Mercado Pago config")
		logger.Error("Failed to initialize Mercado Pago config", err, nil)
		return
	}

	payment = &mercadoPagoPayment{
		config: c,
		log:    logger,
	}

	logger.Info("Mercado Pago payment adapter initialized successfully", nil)
	return
}

func (m *mercadoPagoPayment) CreateCardPayment(ctx context.Context, request bPorts.CardPaymentRequest) (resp *bPorts.CardPaymentResponse, err error) {
	client := payment.NewClient(m.config)

	paymentRequest := payment.Request{
		TransactionAmount: request.TransactionAmount,
		Description:       request.Description,
		Installments:      request.Installments,
		PaymentMethodID:   request.PaymentMethodID,
		Payer: &payment.PayerRequest{
			Email: request.Payer.Email,
		},
	}

	if request.Payer.Identification != nil {
		paymentRequest.Payer.Identification = &payment.IdentificationRequest{
			Type:   request.Payer.Identification.Type,
			Number: request.Payer.Identification.Number,
		}
	}

	if request.ExternalReference != "" {
		paymentRequest.ExternalReference = request.ExternalReference
	}

	if request.NotificationURL != "" {
		paymentRequest.NotificationURL = request.NotificationURL
	}

	cardToken, err := m.GenerateCardToken(ctx, request.BodyCard)
	if err != nil {
		m.log.Error("Failed to generate card token", err, map[string]any{
			"external_reference": request.ExternalReference,
		})
		return nil, fmt.Errorf("failed to generate card token: %w", err)
	}

	paymentRequest.Token = cardToken

	resource, err := client.Create(ctx, paymentRequest)
	if err != nil {
		m.log.Error("Mercado Pago API failed to create payment", err, map[string]any{
			"external_reference": request.ExternalReference,
			"amount":             request.TransactionAmount,
		})
		return nil, fmt.Errorf("failed to create card payment: %w", err)
	}

	m.log.Info("Payment created successfully", map[string]any{
		"payment_id":         resource.ID,
		"status":             resource.Status,
		"external_reference": resource.ExternalReference,
	})

	response := &bPorts.CardPaymentResponse{
		ID:                int64(resource.ID),
		Status:            resource.Status,
		StatusDetail:      resource.StatusDetail,
		TransactionAmount: resource.TransactionAmount,
		Description:       resource.Description,
		ExternalReference: resource.ExternalReference,
		PaymentMethodID:   resource.PaymentMethodID,
		PaymentTypeID:     resource.PaymentTypeID,
		Installments:      resource.Installments,
	}

	if resource.Card.FirstSixDigits != "" {
		response.Card = &bPorts.CardInfo{
			FirstSixDigits: resource.Card.FirstSixDigits,
			LastFourDigits: resource.Card.LastFourDigits,
		}
	}

	return response, nil
}

func (m *mercadoPagoPayment) GenerateCardToken(ctx context.Context, request bPorts.BodyCard) (cardTokenID string, err error) {
	client := cardtoken.NewClient(m.config)

	resource, err := client.Create(context.Background(), cardtoken.Request{
		SiteID:          request.SiteID,
		CardNumber:      request.CardNumber,
		ExpirationYear:  request.ExpirationYear,
		ExpirationMonth: request.ExpirationMonth,
		SecurityCode:    request.SecurityCode,
		Cardholder: &cardtoken.CardholderRequest{
			Identification: &cardtoken.IdentificationRequest{
				Type:   request.Cardholder.Identification.Type,
				Number: request.Cardholder.Identification.Number,
			},
			Name: request.Cardholder.Name,
		},
	})
	if err != nil {
		m.log.Error("Failed to tokenize card data", err, nil)
		return
	}

	cardTokenID = resource.ID

	return
}

func (m *mercadoPagoPayment) GetPayment(ctx context.Context, paymentID int64) (resp *bPorts.CardPaymentResponse, err error) {
	client := payment.NewClient(m.config)

	resource, err := client.Get(ctx, int(paymentID))
	if err != nil {
		m.log.Error("Failed to retrieve payment from Mercado Pago", err, map[string]any{
			"payment_id": paymentID,
		})
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	m.log.Debug("Payment retrieved successfully", map[string]any{
		"payment_id": resource.ID,
		"status":     resource.Status,
	})

	resp = &bPorts.CardPaymentResponse{
		ID:                int64(resource.ID),
		Status:            resource.Status,
		StatusDetail:      resource.StatusDetail,
		TransactionAmount: resource.TransactionAmount,
		Description:       resource.Description,
		ExternalReference: resource.ExternalReference,
		PaymentMethodID:   resource.PaymentMethodID,
		PaymentTypeID:     resource.PaymentTypeID,
		Installments:      resource.Installments,
	}

	if resource.Card.FirstSixDigits != "" {
		resp.Card = &bPorts.CardInfo{
			FirstSixDigits: resource.Card.FirstSixDigits,
			LastFourDigits: resource.Card.LastFourDigits,
		}
	}

	return
}

func (m *mercadoPagoPayment) CreateRefund(ctx context.Context, paymentID int64, amount *float64) (resp *bPorts.RefundResponse, err error) {
	client := refund.NewClient(m.config)

	var resource *refund.Response

	if amount != nil {
		resource, err = client.CreatePartialRefund(ctx, int(paymentID), *amount)
	} else {
		resource, err = client.Create(ctx, int(paymentID))
	}

	if err != nil {
		m.log.Error("Failed to create refund", err, map[string]any{
			"payment_id": paymentID,
			"amount":     amount,
		})
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	refundType := "full"
	if amount != nil {
		refundType = "partial"
	}
	m.log.Info("Refund created successfully", map[string]any{
		"refund_id":   resource.ID,
		"payment_id":  resource.PaymentID,
		"status":      resource.Status,
		"amount":      resource.Amount,
		"refund_type": refundType,
	})

	resp = &bPorts.RefundResponse{
		ID:                   int64(resource.ID),
		PaymentID:            int64(resource.PaymentID),
		Amount:               resource.Amount,
		Status:               resource.Status,
		DateCreated:          resource.DateCreated.Format("2006-01-02T15:04:05.000-07:00"),
		SourceID:             0,
		RefundMode:           resource.RefundMode,
		AdjustmentAmount:     resource.AdjustmentAmount,
		StatusDetail:         resource.Reason,
		UniqueSequenceNumber: resource.UniqueSequenceNumber,
	}

	if resource.Source.ID != "" {
		sourceID := 0
		fmt.Sscanf(resource.Source.ID, "%d", &sourceID)
		resp.SourceID = int64(sourceID)
	}

	return
}

func (m *mercadoPagoPayment) CancelPayment(ctx context.Context, paymentID int64) (err error) {
	client := payment.NewClient(m.config)

	resource, err := client.Cancel(ctx, int(paymentID))
	if err != nil {
		m.log.Error("Failed to cancel payment", err, map[string]any{
			"payment_id": paymentID,
		})
		return fmt.Errorf("failed to cancel payment: %w", err)
	}

	if resource.Status != "cancelled" {
		m.log.Warn("Payment cancellation returned unexpected status", map[string]any{
			"payment_id":     resource.ID,
			"current_status": resource.Status,
		})
		return fmt.Errorf("payment cancellation failed, current status: %s", resource.Status)
	}

	m.log.Info("Payment cancelled successfully", map[string]any{
		"payment_id": resource.ID,
	})

	return
}

func (m *mercadoPagoPayment) SearchPayments(ctx context.Context, filters bPorts.PaymentSearchFilters) (resp *bPorts.PaymentSearchResponse, err error) {
	client := payment.NewClient(m.config)

	searchRequest := payment.SearchRequest{
		Filters: map[string]string{},
	}

	if filters.ExternalReference != "" {
		searchRequest.Filters["external_reference"] = filters.ExternalReference
	}

	if filters.Status != "" {
		searchRequest.Filters["status"] = filters.Status
	}

	if filters.DateFrom != "" {
		searchRequest.Filters["begin_date"] = filters.DateFrom
	}

	if filters.DateTo != "" {
		searchRequest.Filters["end_date"] = filters.DateTo
	}

	if filters.Limit > 0 {
		searchRequest.Limit = filters.Limit
	} else {
		searchRequest.Limit = 30
	}

	if filters.Offset > 0 {
		searchRequest.Offset = filters.Offset
	}

	resource, err := client.Search(ctx, searchRequest)
	if err != nil {
		m.log.Error("Failed to search payments", err, map[string]any{
			"filters": searchRequest.Filters,
			"limit":   searchRequest.Limit,
			"offset":  searchRequest.Offset,
		})
		return nil, fmt.Errorf("failed to search payments: %w", err)
	}

	m.log.Debug("Payment search completed", map[string]any{
		"total_results": resource.Paging.Total,
		"returned":      len(resource.Results),
		"filters":       searchRequest.Filters,
	})

	resp = &bPorts.PaymentSearchResponse{
		Paging: &bPorts.PagingInfo{
			Total:  resource.Paging.Total,
			Limit:  resource.Paging.Limit,
			Offset: resource.Paging.Offset,
		},
		Results: make([]bPorts.CardPaymentResponse, 0, len(resource.Results)),
	}

	for _, paymentResult := range resource.Results {
		paymentResp := bPorts.CardPaymentResponse{
			ID:                int64(paymentResult.ID),
			Status:            paymentResult.Status,
			StatusDetail:      paymentResult.StatusDetail,
			TransactionAmount: paymentResult.TransactionAmount,
			Description:       paymentResult.Description,
			ExternalReference: paymentResult.ExternalReference,
			PaymentMethodID:   paymentResult.PaymentMethodID,
			PaymentTypeID:     paymentResult.PaymentTypeID,
			Installments:      paymentResult.Installments,
		}

		if paymentResult.Card.FirstSixDigits != "" {
			paymentResp.Card = &bPorts.CardInfo{
				FirstSixDigits: paymentResult.Card.FirstSixDigits,
				LastFourDigits: paymentResult.Card.LastFourDigits,
			}
		}

		resp.Results = append(resp.Results, paymentResp)
	}

	return
}

func (m *mercadoPagoPayment) CreatePreference(ctx context.Context, request bPorts.PreferenceRequest) (resp *bPorts.PreferenceResponse, err error) {
	client := preference.NewClient(m.config)

	items := make([]preference.ItemRequest, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, preference.ItemRequest{
			Title:       item.Title,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			CurrencyID:  item.CurrencyID,
			Description: item.Description,
		})
	}

	prefRequest := preference.Request{
		Items:               items,
		ExternalReference:   request.ExternalReference,
		NotificationURL:     request.NotificationURL,
		StatementDescriptor: request.StatementDescriptor,
	}

	if request.BackURLs != nil {
		prefRequest.BackURLs = &preference.BackURLsRequest{
			Success: request.BackURLs.Success,
			Failure: request.BackURLs.Failure,
			Pending: request.BackURLs.Pending,
		}
	}

	if request.AutoReturn != "" {
		prefRequest.AutoReturn = request.AutoReturn
	}

	if request.Payer != nil {
		prefRequest.Payer = &preference.PayerRequest{
			Name:    request.Payer.Name,
			Surname: request.Payer.Surname,
			Email:   request.Payer.Email,
		}
		if request.Payer.Phone != nil {
			prefRequest.Payer.Phone = &preference.PhoneRequest{
				AreaCode: request.Payer.Phone.AreaCode,
				Number:   request.Payer.Phone.Number,
			}
		}
	}

	if request.Expires {
		prefRequest.Expires = request.Expires
	}

	if len(request.Metadata) > 0 {
		prefRequest.Metadata = request.Metadata
	}

	resource, err := client.Create(ctx, prefRequest)
	if err != nil {
		m.log.Error("Failed to create preference (quotation)", err, map[string]any{
			"external_reference": request.ExternalReference,
		})
		return nil, fmt.Errorf("failed to create preference: %w", err)
	}

	m.log.Info("Preference (quotation) created successfully", map[string]any{
		"preference_id":      resource.ID,
		"init_point":         resource.InitPoint,
		"external_reference": resource.ExternalReference,
	})

	resp = &bPorts.PreferenceResponse{
		ID:                resource.ID,
		InitPoint:         resource.InitPoint,
		SandboxInitPoint:  resource.SandboxInitPoint,
		Items:             request.Items,
		ExternalReference: resource.ExternalReference,
		DateCreated:       resource.DateCreated.Format("2006-01-02T15:04:05.000-07:00"),
		CollectorID:       resource.CollectorID,
		ClientID:          resource.ClientID,
		Marketplace:       resource.Marketplace,
		Metadata:          request.Metadata,
	}

	if request.ExpirationDateFrom != nil {
		resp.ExpirationDateFrom = request.ExpirationDateFrom
	}
	if request.ExpirationDateTo != nil {
		resp.ExpirationDateTo = request.ExpirationDateTo
	}
	if request.DateOfExpiration != nil {
		resp.DateOfExpiration = request.DateOfExpiration
	}
	resp.Expires = request.Expires

	return
}

func (m *mercadoPagoPayment) GetPreference(ctx context.Context, preferenceID string) (resource *preference.Response, err error) {
	client := preference.NewClient(m.config)

	resource, err = client.Get(ctx, preferenceID)
	if err != nil {
		m.log.Error("Failed to get preference", err, map[string]any{
			"preference_id": preferenceID,
		})
		return nil, fmt.Errorf("failed to get preference: %w", err)
	}

	m.log.Debug("Preference retrieved successfully", map[string]any{
		"preference_id": resource.ID,
	})

	return
}

func (m *mercadoPagoPayment) UpdatePreference(ctx context.Context, preferenceID string, request bPorts.PreferenceRequest) (resp *bPorts.PreferenceResponse, err error) {
	client := preference.NewClient(m.config)

	items := make([]preference.ItemRequest, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, preference.ItemRequest{
			Title:       item.Title,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			CurrencyID:  item.CurrencyID,
			Description: item.Description,
		})
	}

	prefRequest := preference.Request{
		Items:               items,
		ExternalReference:   request.ExternalReference,
		NotificationURL:     request.NotificationURL,
		StatementDescriptor: request.StatementDescriptor,
	}

	if request.BackURLs != nil {
		prefRequest.BackURLs = &preference.BackURLsRequest{
			Success: request.BackURLs.Success,
			Failure: request.BackURLs.Failure,
			Pending: request.BackURLs.Pending,
		}
	}

	if request.AutoReturn != "" {
		prefRequest.AutoReturn = request.AutoReturn
	}

	if request.Expires {
		prefRequest.Expires = request.Expires

	}

	if len(request.Metadata) > 0 {
		prefRequest.Metadata = request.Metadata
	}

	resource, err := client.Update(ctx, preferenceID, prefRequest)
	if err != nil {
		m.log.Error("Failed to update preference", err, map[string]any{
			"preference_id":      preferenceID,
			"external_reference": request.ExternalReference,
		})
		return nil, fmt.Errorf("failed to update preference: %w", err)
	}

	m.log.Info("Preference updated successfully", map[string]any{
		"preference_id":      resource.ID,
		"external_reference": resource.ExternalReference,
	})

	resp = &bPorts.PreferenceResponse{
		ID:                resource.ID,
		InitPoint:         resource.InitPoint,
		SandboxInitPoint:  resource.SandboxInitPoint,
		Items:             request.Items,
		ExternalReference: resource.ExternalReference,
		DateCreated:       resource.DateCreated.Format("2006-01-02T15:04:05.000-07:00"),
		CollectorID:       int64(resource.CollectorID),
		ClientID:          resource.ClientID,
		Marketplace:       resource.Marketplace,
		Metadata:          request.Metadata,
	}

	if request.ExpirationDateFrom != nil {
		resp.ExpirationDateFrom = request.ExpirationDateFrom
	}
	if request.ExpirationDateTo != nil {
		resp.ExpirationDateTo = request.ExpirationDateTo
	}
	if request.DateOfExpiration != nil {
		resp.DateOfExpiration = request.DateOfExpiration
	}
	resp.Expires = request.Expires

	return
}
