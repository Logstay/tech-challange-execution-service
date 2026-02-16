package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type ListServiceDefinitionHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler ListServiceDefinitionHandler) Handle(ctx context.Context,
	filter ports.Filter) (result []ports.GetServiceDefinition, total int64, err error) {

	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-list-service-definition")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-list-service-definition-params", tracing.Attributes{
		"filter": filter,
	})

	result, total, err = handler.Repository.List(ctx, filter)
	if err != nil {
		span.AddEvent("service-list-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("service-list-service-definition-success", nil)

	return

}
