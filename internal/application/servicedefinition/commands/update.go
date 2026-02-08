package command

import (
	"context"

	model "github.com/Logstay/tech-challange-wherehouse-service/internal/domain/servicedefinition"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/servicedefinition"
	httperr "github.com/Logstay/tech-challange-wherehouse-service/pkg/errors"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type UpdateServiceDefinitionHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler UpdateServiceDefinitionHandler) Handle(ctx context.Context,
	cmd ports.UpdateServiceDefinition) (err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-update-service-definition")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-update-service-definition-params", tracing.Attributes{
		"cmd": cmd,
	})

	resp, err := handler.Repository.Get(ctx, cmd.UID.String)
	if err != nil {
		span.AddEvent("service-update-consumable-exist-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if resp.UID.IsZero() {
		return httperr.NotFoundError("serviço não encontrado")
	}

	serviceDefinition, err := handler.serviceDefinition(cmd)
	if err != nil {
		span.AddEvent("service-update-service-definition-instance-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = handler.Repository.Update(ctx, cmd.UID.String, serviceDefinition)
	if err != nil {
		span.AddEvent("service-update-service-definition-repository-update-error", tracing.Attributes{
			"err": err,
		})
		return err
	}

	span.AddEvent("service-update-service-definition-success", nil)

	return nil

}

func (handler UpdateServiceDefinitionHandler) serviceDefinition(
	cmd ports.UpdateServiceDefinition) (result model.ServiceDefinition, err error) {
	serviceDefinition := model.NewServiceDefinitionFactory()

	if err = serviceDefinition.Price(cmd.Price); err != nil {
		return
	}

	if err = serviceDefinition.Name(cmd.Name); err != nil {
		return
	}

	result = serviceDefinition.Build()

	return
}
