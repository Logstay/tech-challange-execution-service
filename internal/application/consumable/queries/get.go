package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
	httperr "github.com/Logstay/tech-challange-wherehouse-service/pkg/errors"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type GetConsumableHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler GetConsumableHandler) Handle(ctx context.Context, uid string) (result ports.GetConsumable, err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-get-consumable")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-get-consumable-params", tracing.Attributes{
		"uid": uid,
	})

	result, err = handler.Repository.Get(ctx, uid)
	if err != nil {
		span.AddEvent("service-get-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if result.UID.IsZero() {
		err = httperr.NotFoundError("peça/insumo não encontrado")
		return
	}

	span.AddEvent("service-get-consumable-success", nil)

	return

}
