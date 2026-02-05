package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/billing"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/logger"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type GetQuotationHandler struct {
	BillingAdapter ports.MercadoPagoBilling
	Logger         logger.Logger
	Tracer         *tracing.Tracer
}

func (handler GetQuotationHandler) Handle(ctx context.Context, id string) (result *ports.PreferenceResponse, err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-get-quotation")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-get-quotation-params", tracing.Attributes{
		"id": id,
	})

	resp, err := handler.BillingAdapter.GetPreference(ctx, id)
	if err != nil {
		span.AddEvent("service-get-quotation-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	items := make([]ports.PreferenceItemRequest, 0, len(resp.Items))
	for _, item := range resp.Items {
		items = append(items, ports.PreferenceItemRequest{
			Title:       item.Title,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			CurrencyID:  item.CurrencyID,
			Description: item.Description,
		})
	}

	result = &ports.PreferenceResponse{
		ID:                resp.ID,
		InitPoint:         resp.InitPoint,
		SandboxInitPoint:  resp.SandboxInitPoint,
		Items:             items,
		Status:            resp.AutoReturn,
		ExternalReference: resp.ExternalReference,
		DateCreated:       resp.DateCreated.Format("2006-01-02T15:04:05.000-07:00"),
		CollectorID:       resp.CollectorID,
		ClientID:          resp.ClientID,
		Marketplace:       resp.Marketplace,
	}

	span.AddEvent("service-get-quotation-success", nil)

	return

}
