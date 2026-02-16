package consumable

import (
	"context"

	modelconsumable "github.com/Logstay/tech-challange-execution-service/internal/domain/consumable"
	ports "github.com/Logstay/tech-challange-execution-service/internal/ports/consumable"
	inventoryports "github.com/Logstay/tech-challange-execution-service/internal/ports/inventory"
	"github.com/Logstay/tech-challange-execution-service/pkg/null"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type repository struct {
	dao          ports.DAO
	tracer       *tracing.Tracer
	inventoryDAO inventoryports.DAO
}

func NewRepository(
	dao ports.DAO,
	inventoryDAO inventoryports.DAO,
	tracingProvider *tracing.Provider) ports.Repository {

	return &repository{
		dao:          dao,
		inventoryDAO: inventoryDAO,
		tracer:       tracingProvider.NewTracer("repository-consumable"),
	}
}

func (r *repository) Save(ctx context.Context, in modelconsumable.Consumable) (err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-save-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-save-consumable-params", tracing.Attributes{
		"in": in,
	})

	id, err := r.dao.Save(ctx, ports.ConsumableDAO{
		Name:        in.Name(),
		Price:       in.Price(),
		Code:        in.Code(),
		Description: in.Description(),
		CreatedAt:   in.CreatedAt(),
		UpdatedAt:   in.UpdatedAt(),
	})
	if err != nil {
		span.AddEvent("repository-save-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = r.inventoryDAO.Save(ctx, inventoryports.InventoryDAO{
		ConsumableID: null.IntFrom(id),
		Quantity:     in.Quantity(),
	})
	if err != nil {
		span.AddEvent("repository-save-inventory-save-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-save-consumable-success", nil)
	return
}

func (r *repository) Get(ctx context.Context, uid string) (result ports.GetConsumable, err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-get-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-get-consumable-params", tracing.Attributes{
		"uid": uid,
	})

	resp, err := r.dao.Get(ctx, uid)
	if err != nil {
		span.AddEvent("repository-get-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	inventory, err := r.inventoryDAO.Get(ctx, resp.ID.Int64)
	if err != nil {
		span.AddEvent("repository-get-consumable-inventory-get-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	result = ports.GetConsumable{
		UID:         resp.UID,
		Name:        resp.Name,
		Description: resp.Description,
		Price:       resp.Price,
		Code:        resp.Code,
		CreatedAt:   resp.CreatedAt,
		Quantity:    inventory.Quantity,
		UpdatedAt:   resp.UpdatedAt,
	}

	span.AddEvent("repository-get-consumable-success", nil)
	return
}

func (r *repository) List(ctx context.Context, filter ports.Filter) (result []ports.GetConsumable, total int64, err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-list-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-list-consumable-params", tracing.Attributes{
		"filter": filter,
	})

	resp, err := r.dao.List(ctx, filter)
	if err != nil {
		span.AddEvent("repository-list-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	for _, v := range resp {
		var inventory inventoryports.InventoryDAO

		inventory, err = r.inventoryDAO.Get(ctx, v.ID.Int64)
		if err != nil {
			span.AddEvent("repository-list-consumable-inventory-get-error", tracing.Attributes{
				"err": err,
			})
			return
		}

		result = append(result, ports.GetConsumable{
			UID:         v.UID,
			Name:        v.Name,
			Description: v.Description,
			Code:        v.Code,
			Price:       v.Price,
			CreatedAt:   v.CreatedAt,
			Quantity:    inventory.Quantity,
			UpdatedAt:   v.UpdatedAt,
		})
	}

	total, err = r.dao.TotalList(ctx, filter)
	if err != nil {
		span.AddEvent("repository-list-consumable-total-list-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-list-consumable-success", nil)
	return
}

func (r *repository) Update(ctx context.Context, uid string, in modelconsumable.Consumable) (err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-update-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-update-consumable-params", tracing.Attributes{
		"uid": uid,
		"in":  in,
	})

	err = r.dao.Update(ctx, ports.ConsumableDAO{
		UID:         null.StringFrom(uid),
		Name:        in.Name(),
		Price:       in.Price(),
		Description: in.Description(),
		CreatedAt:   in.CreatedAt(),
		UpdatedAt:   in.UpdatedAt(),
	})
	if err != nil {
		span.AddEvent("repository-update-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-update-consumable-success", nil)
	return
}

func (r *repository) Delete(ctx context.Context, uid string) (err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-delete-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-delete-consumable-params", tracing.Attributes{
		"uid": uid,
	})

	resp, err := r.dao.Get(ctx, uid)
	if err != nil {
		span.AddEvent("repository-get-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = r.inventoryDAO.Delete(ctx, resp.ID.Int64)
	if err != nil {
		span.AddEvent("repository-delete-consumable-inventory-delete-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = r.dao.Delete(ctx, resp.UID.String)
	if err != nil {
		span.AddEvent("repository-delete-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-delete-consumable-success", nil)
	return
}

func (r *repository) Exist(ctx context.Context, name string) (result bool, err error) {
	var (
		span *tracing.Span
	)

	ctx, span = r.tracer.NewSpan(ctx, "repository-exist-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("repository-exist-consumable-params", tracing.Attributes{
		"name": name,
	})

	result, err = r.dao.Exist(ctx, name)
	if err != nil {
		span.AddEvent("repository-exist-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("repository-exist-consumable-success", nil)
	return
}
