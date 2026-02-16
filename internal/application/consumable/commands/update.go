package command

import (
	"context"

	model "github.com/Logstay/tech-challange-execution-service/internal/domain/consumable"
	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type UpdateConsumableHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler UpdateConsumableHandler) Handle(ctx context.Context, cmd ports.UpdateConsumable) (err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-update-consumable")
	defer func() {
		span.End()
	}()

	resp, err := handler.Repository.Get(ctx, cmd.UID.String)
	if err != nil {
		span.AddEvent("service-delete-consumable-exist-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if resp.UID.IsZero() {
		return httperr.NotFoundError("peça/insumo não encontrado")
	}

	consumable, err := handler.consumable(cmd)
	if err != nil {
		span.AddEvent("service-update-consumable-instance-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = handler.Repository.Update(ctx, cmd.UID.String, consumable)
	if err != nil {
		span.AddEvent("service-update-consumable-update-error", tracing.Attributes{
			"err": err,
		})
		return err
	}

	span.AddEvent("service-update-consumable-success", nil)

	return nil

}

func (handler UpdateConsumableHandler) consumable(
	cmd ports.UpdateConsumable) (result model.Consumable, err error) {
	consumable := model.NewConsumableFactory()

	if err = consumable.Price(cmd.Price); err != nil {
		return
	}

	if err = consumable.Name(cmd.Name); err != nil {
		return
	}

	result = consumable.Build()

	return
}
