package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type ListConsumableHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler ListConsumableHandler) Handle(ctx context.Context,
	filter ports.Filter) (result []ports.GetConsumable, total int64, err error) {

	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-list-consumable")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-list-consumable-params", tracing.Attributes{
		"filter": filter,
	})

	result, total, err = handler.Repository.List(ctx, filter)
	if err != nil {
		span.AddEvent("service-list-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("service-list-consumable-success", nil)

	return

}
