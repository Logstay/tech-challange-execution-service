package servicedefinition

import (
	"context"

	modelservicedefinition "github.com/Logstay/tech-challange-wherehouse-service/internal/domain/servicedefinition"
	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/servicedefinition"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type repository struct {
	dao    ports.DAO
	tracer *tracing.Tracer
}

func NewRepository(
	dao ports.DAO,
	tracingProvider *tracing.Provider) ports.Repository {

	return &repository{
		dao:    dao,
		tracer: tracingProvider.NewTracer("repository-service-definition"),
	}
}

func (r *repository) Save(ctx context.Context, in modelservicedefinition.ServiceDefinition) (err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-save-service-definition")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-save-service-definition-params", tracing.Attributes{
		"in": in,
	})

	_, err = r.dao.Save(ctx, ports.ServiceDefinitionDAO{
		Name:        in.Name(),
		Price:       in.Price(),
		Description: in.Description(),
		CreatedAt:   in.CreatedAt(),
		UpdatedAt:   in.UpdatedAt(),
	})
	if err != nil {
		span.AddEvent("repository-save-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-save-service-definition-success", nil)
	return
}

func (r *repository) Get(ctx context.Context, uid string) (result ports.GetServiceDefinition, err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-get-service-definition")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-get-service-definition-params", tracing.Attributes{
		"uid": uid,
	})

	resp, err := r.dao.Get(ctx, uid)
	if err != nil {
		span.AddEvent("repository-get-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	result = ports.GetServiceDefinition(resp)

	span.AddEvent("repository-get-service-definition-success", nil)
	return
}

func (r *repository) List(ctx context.Context, filter ports.Filter) (result []ports.GetServiceDefinition, total int64, err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-list-service-definition")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-list-service-definition-params", tracing.Attributes{
		"filter": filter,
	})

	resp, err := r.dao.List(ctx, filter)
	if err != nil {
		span.AddEvent("repository-list-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	for _, v := range resp {
		result = append(result, ports.GetServiceDefinition(v))
	}

	total, err = r.dao.TotalList(ctx, filter)
	if err != nil {
		span.AddEvent("repository-list-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-list-service-definition-success", nil)
	return
}

func (r *repository) Update(ctx context.Context, uid string, in modelservicedefinition.ServiceDefinition) (err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-update-service-definition")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-update-service-definition-params", tracing.Attributes{
		"uid": uid,
		"in":  in,
	})

	err = r.dao.Update(ctx, ports.ServiceDefinitionDAO{
		UID:         null.StringFrom(uid),
		Name:        in.Name(),
		Price:       in.Price(),
		Description: in.Description(),
		CreatedAt:   in.CreatedAt(),
		UpdatedAt:   in.UpdatedAt(),
	})
	if err != nil {
		span.AddEvent("repository-update-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-update-service-definition-success", nil)
	return
}

func (r *repository) Delete(ctx context.Context, uid string) (err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-delete-service-definition")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-delete-service-definition-params", tracing.Attributes{
		"uid": uid,
	})

	err = r.dao.Delete(ctx, uid)
	if err != nil {
		span.AddEvent("repository-delete-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-delete-service-definition-success", nil)
	return
}

func (r *repository) Exist(ctx context.Context, name string) (result bool, err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-exist-service-definition")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-exist-service-definition-params", tracing.Attributes{
		"name": name,
	})

	result, err = r.dao.Exist(ctx, name)
	if err != nil {
		span.AddEvent("repository-exist-service-definition-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-exist-service-definition-success", nil)
	return
}
