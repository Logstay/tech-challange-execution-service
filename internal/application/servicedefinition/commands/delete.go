package command

import (
	"context"

	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
	httperr "github.com/Logstay/tech-challange-execution-service/pkg/errors"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type DeleteServiceDefinitionHandler struct {
	Repository ports.Repository
	Tracer     *tracing.Tracer
}

func (handler DeleteServiceDefinitionHandler) Handle(ctx context.Context, uid string) (err error) {
	var span *tracing.Span

	ctx, span = handler.Tracer.NewSpan(ctx, "service-delete-service-definition")
	defer func() {
		span.End()
	}()

	span.AddEvent("service-delete-service-definition-params", tracing.Attributes{
		"uid": uid,
	})

	resp, err := handler.Repository.Get(ctx, uid)
	if err != nil {
		span.AddEvent("service-delete-consumable-get-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	if resp.UID.IsZero() {
		return httperr.NotFoundError("serviço não encontrado")
	}

	err = handler.Repository.Delete(ctx, uid)
	if err != nil {
		span.AddEvent("service-delete-service-definition-delete-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("service-delete-service-definition-success", nil)
	return
}
