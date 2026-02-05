package billing

import (
	"context"

	qries "github.com/Logstay/tech-challange-wherehouse-service/internal/application/billing/queries"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/billing"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/logger"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type Service struct {
	// queries
	getQuotationQuery *qries.GetQuotationHandler

	d Dependencies
}

type Dependencies struct {
	BillingAdapter ports.MercadoPagoBilling
	Tracer         *tracing.Tracer
	Logger         logger.Logger
}

func NewService(d Dependencies) (s *Service) {
	s = &Service{
		d: d,
	}

	s.buildQueries()

	return
}

// queries
func (s *Service) GetQuotationByID(ctx context.Context, id string) (result *ports.PreferenceResponse, err error) {
	return s.getQuotationQuery.Handle(ctx, id)
}

func (s *Service) buildQueries() {
	s.buildGetQuotation()
}

func (s *Service) buildGetQuotation() {
	s.getQuotationQuery = &qries.GetQuotationHandler{
		BillingAdapter: s.d.BillingAdapter,
		Tracer:         s.d.Tracer,
		Logger:         s.d.Logger,
	}
}
