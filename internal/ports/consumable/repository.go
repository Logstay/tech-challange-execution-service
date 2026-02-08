package ports

import (
	"context"

	modelconsumable "github.com/Logstay/tech-challange-wherehouse-service/internal/domain/consumable"
)

type Repository interface {
	Save(ctx context.Context, in modelconsumable.Consumable) (err error)
	Get(ctx context.Context, uid string) (result GetConsumable, err error)
	List(ctx context.Context, filter Filter) (result []GetConsumable, total int64, err error)
	Update(ctx context.Context, uid string, in modelconsumable.Consumable) (err error)
	Delete(ctx context.Context, uid string) (err error)
	Exist(ctx context.Context, name string) (result bool, err error)
}
