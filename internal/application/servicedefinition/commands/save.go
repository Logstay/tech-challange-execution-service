package command

import (
	"context"

	modelserviceDefinition "github.com/Logstay/tech-challange-execution-service/internal/domain/servicedefinition"
	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type SaveServiceDefinitionHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler SaveServiceDefinitionHandler) Handle(ctx context.Context, cmd ports.SaveServiceDefinition) (err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-save-service-definition")
	defer func() {
		span.End()
	}()

	exist, err := handler.Repository.Exist(ctx, cmd.Name.String)
	if err != nil {
		span.AddEvent("service-save-service-definition-repository-exist-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if exist {
		span.AddEvent("service-save-service-definition-repository-exist-error", tracing.Attributes{
			"err": "serviceDefinition already exist",
		})
		return httperr.ExistError("existe um serviço cadastrado com esse nome")
	}

	serviceDefinition, err := handler.serviceDefinition(cmd)
	if err != nil {
		span.AddEvent("service-save-service-definition-instance-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = handler.Repository.Save(ctx, serviceDefinition)
	if err != nil {
		span.AddEvent("service-save-service-definition-repository-save-error", tracing.Attributes{
			"err": err,
		})
		return err
	}

	span.AddEvent("service-save-service-definition-success", nil)

	return nil

}

func (handler SaveServiceDefinitionHandler) serviceDefinition(
	cmd ports.SaveServiceDefinition) (result modelserviceDefinition.ServiceDefinition, err error) {
	serviceDefinition := modelserviceDefinition.NewServiceDefinitionFactory()

	if err = serviceDefinition.Price(cmd.Price); err != nil {
		return
	}

	if err = serviceDefinition.Name(cmd.Name); err != nil {
		return
	}

	if err = serviceDefinition.Description(cmd.Description); err != nil {
		return
	}

	result = serviceDefinition.Build()

	return
}
