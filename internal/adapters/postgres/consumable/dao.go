package consumable

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"

	ports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
	httperr "github.com/Logstay/tech-challange-wherehouse-service/pkg/errors"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
	database "github.com/Logstay/tech-challange-wherehouse-service/pkg/postgres"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/sqlbuilder"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
	"github.com/jmoiron/sqlx"
)

const (
	regex = "[^a-zA-Z0-9]+"
)

type DAO struct {
	db     database.UnitOfWork
	tracer *tracing.Tracer
}

type consumable struct {
	ID          null.Int      `db:"id"`
	UID         null.String   `db:"uid"`
	Name        null.String   `db:"name"`
	Code        null.String   `db:"code"`
	Description null.String   `db:"description"`
	Price       null.Float    `db:"price"`
	CreatedAt   null.DateTime `db:"created_at"`
	UpdatedAt   null.DateTime `db:"updated_at"`
}

func NewDAO(db database.UnitOfWork, tracingProvider *tracing.Provider) *DAO {
	return &DAO{
		db:     db,
		tracer: tracingProvider.NewTracer("dao-consumable"),
	}
}

func (d *DAO) Save(ctx context.Context, input ports.ConsumableDAO) (id int64, err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-save-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-save-consumable-params", tracing.Attributes{
		"input": input,
	})

	res, err := d.db.NamedQuery(qrySaveConsumable, d.modelTo(input))
	if err != nil {
		span.AddEvent("dao-save-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	id, err = database.CloseDBTransactionReturningID(res, err)

	span.AddEvent("dao-save-consumable-success", nil)

	return
}

func (d *DAO) Get(ctx context.Context, uid string) (result ports.ConsumableDAO, err error) {
	var (
		span *tracing.Span
		resp consumable
	)
	_, span = d.tracer.NewSpan(ctx, "dao-get-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-get-consumable-params", tracing.Attributes{
		"uid": uid,
	})

	err = d.db.Get(&resp, qryGetConsumable, uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		span.AddEvent("dao-get-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	result = d.modelFrom(resp)

	span.AddEvent("dao-get-consumable-success", nil)

	return
}

func (d *DAO) GetByID(ctx context.Context, id int64) (result ports.ConsumableDAO, err error) {
	var (
		span *tracing.Span
		resp consumable
	)
	_, span = d.tracer.NewSpan(ctx, "dao-get-by-id-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-get-by-id-consumable-params", tracing.Attributes{
		"id": id,
	})

	err = d.db.Get(&resp, qryGetConsumableByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		span.AddEvent("dao-get-by-id-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	result = d.modelFrom(resp)

	span.AddEvent("dao-get-by-id-consumable-success", nil)

	return
}

func (d *DAO) List(ctx context.Context, filter ports.Filter) (result []ports.ConsumableDAO, err error) {
	var (
		span *tracing.Span
	)

	resp := []consumable{}

	_, span = d.tracer.NewSpan(ctx, "dao-list-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-list-consumable-params", tracing.Attributes{
		"filter": filter,
	})

	sql, args := sqlbuilder.Select(qryListConsumable).Where(func(args sqlbuilder.Value) {

		if filter.Name.String != "" {
			reg := regexp.MustCompile(regex)
			texto := reg.ReplaceAllString(filter.Name.String, "")
			texto = "%" + strings.ToUpper(texto) + "%"
			args.Add("AND UPPER(regexp_replace(c.name, '[^a-zA-Z0-9]+', '', 'g')) LIKE ?", texto)
		}
	}).
		OrderBy("ORDER BY c.created_at DESC, c.name").
		SetPaginate(filter.Page.Int64*filter.Limit.Int64, filter.Limit.Int64).
		Build()

	err = d.db.Select(&resp, sqlx.Rebind(sqlx.DOLLAR, sql), args...)
	if err != nil {
		err = httperr.InternalServerError("internal error")
		return
	}

	for _, v := range resp {
		result = append(result, d.modelFrom(v))
	}

	return
}

func (d *DAO) TotalList(ctx context.Context, filter ports.Filter) (total int64, err error) {
	var (
		span *tracing.Span
	)

	_, span = d.tracer.NewSpan(ctx, "dao-total-list-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-total-list-consumable-params", tracing.Attributes{
		"filter": filter,
	})

	sql, args := sqlbuilder.Select(qryTotalListConsumable).Where(func(args sqlbuilder.Value) {

		if filter.Name.String != "" {
			reg := regexp.MustCompile(regex)
			texto := reg.ReplaceAllString(filter.Name.String, "")
			texto = "%" + strings.ToUpper(texto) + "%"
			args.Add("AND UPPER(regexp_replace(c.name, '[^a-zA-Z0-9]+', '', 'g')) LIKE ?", texto)
		}
	}).
		Build()

	err = d.db.Get(&total, sqlx.Rebind(sqlx.DOLLAR, sql), args...)
	if err != nil {
		err = httperr.InternalServerError("internal error")
		return
	}

	span.AddEvent("dao-total-list-consumable-success", nil)

	return
}

func (d *DAO) Update(ctx context.Context, input ports.ConsumableDAO) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-update-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-update-consumable-params", tracing.Attributes{
		"input": input,
	})

	_, err = d.db.MustNamedExec(qryUpdateConsumable, d.modelTo(input)).RowsAffected()
	if err != nil {
		span.AddEvent("dao-update-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("dao-update-consumable-success", nil)

	return
}

func (d *DAO) Delete(ctx context.Context, uid string) (err error) {
	var span *tracing.Span

	_, span = d.tracer.NewSpan(ctx, "dao-delete-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-delete-consumable-params", tracing.Attributes{
		"uid": uid,
	})

	_, err = d.db.MustExec(qryDeleteConsumable, uid).RowsAffected()
	if err != nil {
		span.AddEvent("dao-delete-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	span.AddEvent("dao-delete-consumable-success", nil)

	return
}

func (d *DAO) Exist(ctx context.Context, name string) (exist bool, err error) {
	var (
		span  *tracing.Span
		count int64
	)
	_, span = d.tracer.NewSpan(ctx, "dao-exist-consumable")
	defer func() {
		span.End(err)
	}()

	span.AddEvent("dao-exist-consumable-params", tracing.Attributes{
		"name": name,
	})

	err = d.db.Get(&count, qryExistConsumable, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = nil
			return
		}
		span.AddEvent("dao-exist-consumable-error", tracing.Attributes{
			"err": err,
		})
		return
	}

	exist = count != 0

	span.AddEvent("dao-exist-consumable-success", nil)

	return
}

func (d *DAO) modelTo(model ports.ConsumableDAO) consumable {
	return consumable{
		ID:          model.ID,
		UID:         model.UID,
		Name:        model.Name,
		Description: model.Description,
		Code:        model.Code,
		Price:       model.Price,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func (d *DAO) modelFrom(model consumable) ports.ConsumableDAO {
	return ports.ConsumableDAO{
		ID:          model.ID,
		UID:         model.UID,
		Name:        model.Name,
		Description: model.Description,
		Code:        model.Code,
		Price:       model.Price,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
