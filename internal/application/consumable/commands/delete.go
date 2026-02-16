package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type DeleteConsumableHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler DeleteConsumableHandler) Handle(ctx context.Context, uid string) (err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-delete-consumable")
	defer func() {
		span.End()
	}()

	resp, err := handler.Repository.Get(ctx, uid)
	if err != nil {
		span.AddEvent("service-delete-consumable-exist-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if resp.UID.IsZero() {
		return httperr.NotFoundError("peça/insumo não encontrado")
	}

	err = handler.Repository.Delete(ctx, uid)
	if err != nil {
		span.AddEvent("service-delete-consumable-delete-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("service-delete-consumable-success", nil)
	return
}
