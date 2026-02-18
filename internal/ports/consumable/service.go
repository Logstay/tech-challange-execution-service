package ports

import (
	"context"

	"github.com/Logstay/tech-challange-execution-service/pkg/null"
)

type Consumable interface {
	SaveConsumable(ctx context.Context, cmd SaveConsumable) (err error)
	UpdateConsumable(ctx context.Context, cmd UpdateConsumable) (err error)
	GetConsumableByUID(ctx context.Context, uid string) (result GetConsumable, err error)
	ListConsumable(ctx context.Context, filter Filter) (result []GetConsumable, total int64, err error)
	DeleteConsumable(ctx context.Context, uid string) (err error)
}

type SaveConsumable struct {
	Name        null.String
	Description null.String
	Code        null.String
	Quantity    null.Int
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}

type GetConsumable struct {
	UID         null.String
	Name        null.String
	Description null.String
	Code        null.String
	Quantity    null.Int
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}

type UpdateConsumable struct {
	UID         null.String
	Name        null.String
	Description null.String
	Code        null.String
	Quantity    null.Int
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}
