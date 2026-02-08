package command

import (
	"context"

	modelconsumable "github.com/Logstay/tech-challange-wherehouse-service/internal/domain/consumable"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
	httperr "github.com/Logstay/tech-challange-wherehouse-service/pkg/errors"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type SaveConsumableHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler SaveConsumableHandler) Handle(ctx context.Context, cmd ports.SaveConsumable) (err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-save-consumable")
	defer func() {
		span.End()
	}()

	exist, err := handler.Repository.Exist(ctx, cmd.Name.String)
	if err != nil {
		span.AddEvent("repository-save-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if exist {
		span.AddEvent("repository-save-consumable-error", tracing.Attributes{
			"err": "consumable already exist",
		})
		return httperr.ConflictError("existe peça/insumo cadastrado com esse nome")
	}

	consumable, err := handler.consumable(cmd)
	if err != nil {
		span.AddEvent("service-save-consumable-instance-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = handler.Repository.Save(ctx, consumable)
	if err != nil {
		span.AddEvent("service-save-consumable-repository-save-error", tracing.Attributes{
			"err": err,
		})
		return err
	}

	span.AddEvent("service-save-consumable-success", nil)

	return nil

}

func (handler SaveConsumableHandler) consumable(
	cmd ports.SaveConsumable) (result modelconsumable.Consumable, err error) {
	consumable := modelconsumable.NewConsumableFactory()

	if err = consumable.Price(cmd.Price); err != nil {
		return
	}

	if err = consumable.Name(cmd.Name); err != nil {
		return
	}

	if err = consumable.Description(cmd.Description); err != nil {
		return
	}

	if err = consumable.Quantity(cmd.Quantity); err != nil {
		return
	}
	result = consumable.Build()

	return
}
