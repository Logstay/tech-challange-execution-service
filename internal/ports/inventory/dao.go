package ports

import (
	"context"

	"github.com/Logstay/tech-challange-wherehouse-service/pkg/null"
)

type DAO interface {
	Save(ctx context.Context, input InventoryDAO) (err error)
	Get(ctx context.Context, consumableID int64) (result InventoryDAO, err error)
	Update(ctx context.Context, input InventoryDAO) (err error)
	Delete(ctx context.Context, consumableID int64) (err error)
	Decrease(ctx context.Context, input InventoryDAO) (err error)
	Increase(ctx context.Context, input InventoryDAO) (err error)
}

type InventoryDAO struct {
	ID           null.Int
	UID          null.String
	Quantity     null.Int
	ConsumableID null.Int
	UpdatedAt    null.DateTime
}

type Filter struct {
	Name  null.String
	Page  null.Int
	Limit null.Int
}
