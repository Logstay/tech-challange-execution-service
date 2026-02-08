package ports

import (
	"context"

	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
)

type DAO interface {
	Save(ctx context.Context, input ConsumableDAO) (id int64, err error)
	Get(ctx context.Context, uid string) (result ConsumableDAO, err error)
	List(ctx context.Context, filter Filter) (result []ConsumableDAO, err error)
	Update(ctx context.Context, input ConsumableDAO) (err error)
	Delete(ctx context.Context, uid string) (err error)
	Exist(ctx context.Context, name string) (exist bool, err error)
	GetByID(ctx context.Context, id int64) (result ConsumableDAO, err error)
	TotalList(ctx context.Context, filter Filter) (total int64, err error)
}

type ConsumableDAO struct {
	ID          null.Int
	UID         null.String
	Name        null.String
	Code        null.String
	Description null.String
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}

type Filter struct {
	Name  null.String
	Page  null.Int
	Limit null.Int
}
