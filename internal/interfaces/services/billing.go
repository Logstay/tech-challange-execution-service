package services

import (
	"context"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/billing"
)

type BillingService interface {
	// queries
	GetQuotationByID(ctx context.Context, id string) (result *ports.PreferenceResponse, err error)
}
