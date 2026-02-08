package inventory

import (
	"context"
	"database/sql"
	"errors"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/inventory"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	database "github.com/Logstay/tech-challange-wherehouse-service/pkg/postgres"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type DAO struct {
	db     database.UnitOfWork
	tracer *tracing.Tracer
}

type inventory struct {
	ID           null.Int      `db:"id"`
	UID          null.String   `db:"uid"`
	Quantity     null.Int      `db:"quantity"`
	ConsumableID null.Int      `db:"consumable_id"`
	UpdatedAt    null.DateTime `db:"updated_at"`
}

func NewDAO(db database.UnitOfWork, tracingProvider *tracing.Provider) *DAO {
	return &DAO{
		db:     db,
		tracer: tracingProvider.NewTracer("dao-inventory"),
	}
}

func (d *DAO) Save(ctx context.Context, input ports.InventoryDAO) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-save-inventory")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-save-inventory-params", tracing.Attributes{
		"input": input,
	})

	res, err := d.db.NamedQuery(qrySaveInventory, d.modelTo(input))
	if err != nil {
		span.AddEvent("dao-save-inventory-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	err = database.CloseDBTransaction(res, err)

	span.AddEvent("dao-save-inventory-success", nil)

	return
}

func (d *DAO) Get(ctx context.Context, consumableID int64) (result ports.InventoryDAO, err error) {
	var (
		span *tracing.Span
		resp inventory
	)
	_, span = d.tracer.NewSpan(ctx, "dao-get-inventory")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-get-inventory-params", tracing.Attributes{
		"consumableID": consumableID,
	})

	err = d.db.Get(&resp, qryGetInventory, consumableID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		span.AddEvent("dao-get-inventory-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	result = d.modelFrom(resp)

	span.AddEvent("dao-get-inventory-success", nil)

	return
}

func (d *DAO) Update(ctx context.Context, input ports.InventoryDAO) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-update-inventory")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-update-inventory-params", tracing.Attributes{
		"input": input,
	})

	_, err = d.db.MustNamedExec(qryUpdateInventory, d.modelTo(input)).RowsAffected()
	if err != nil {
		span.AddEvent("dao-update-inventory-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("dao-update-inventory-success", nil)

	return
}

func (d *DAO) Delete(ctx context.Context, consumableID int64) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-delete-inventory")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-delete-inventory-params", tracing.Attributes{
		"consumableID": consumableID,
	})

	_, err = d.db.MustExec(qryDeleteInventory, consumableID).RowsAffected()
	if err != nil {
		span.AddEvent("dao-delete-inventory-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("dao-delete-inventory-success", nil)

	return
}

func (d *DAO) Decrease(ctx context.Context, input ports.InventoryDAO) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-decrease-inventory")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-decrease-inventory-params", tracing.Attributes{
		"input": input,
	})

	_, err = d.db.MustNamedExec(qryDecreseQuantityInventory, d.modelTo(input)).RowsAffected()
	if err != nil {
		span.AddEvent("dao-decrease-inventory-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("dao-decrease-inventory-success", nil)

	return
}

func (d *DAO) Increase(ctx context.Context, input ports.InventoryDAO) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-increase-inventory")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-increase-inventory-params", tracing.Attributes{
		"input": input,
	})

	_, err = d.db.MustNamedExec(qryIncreaseQuantityInventory, d.modelTo(input)).RowsAffected()
	if err != nil {
		span.AddEvent("dao-increase-inventory-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("dao-increase-inventory-success", nil)

	return
}

func (d *DAO) modelTo(model ports.InventoryDAO) inventory {
	return inventory{
		ID:           model.ID,
		UID:          model.UID,
		Quantity:     model.Quantity,
		ConsumableID: model.ConsumableID,
		UpdatedAt:    model.UpdatedAt,
	}
}

func (d *DAO) modelFrom(model inventory) ports.InventoryDAO {
	return ports.InventoryDAO{
		ID:           model.ID,
		UID:          model.UID,
		Quantity:     model.Quantity,
		ConsumableID: model.ConsumableID,
		UpdatedAt:    model.UpdatedAt,
	}
}
