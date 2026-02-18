package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type GetServiceDefinitionHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler GetServiceDefinitionHandler) Handle(ctx context.Context, uid string) (result ports.GetServiceDefinition, err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-get-service-definition")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-get-service-definition-params", tracing.Attributes{
		"uid": uid,
	})

	result, err = handler.Repository.Get(ctx, uid)
	if err != nil {
		span.AddEvent("service-get-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if result.UID.IsZero() {
		err = httperr.NotFoundError("serviço não encontrado")
		return
	}

	span.AddEvent("service-get-service-definition-success", nil)

	return

}
