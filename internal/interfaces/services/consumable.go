package services

import (
	"context"

	cports "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/consumable"
)

type ConsumableService interface {
	SaveConsumable(ctx context.Context, cmd cports.SaveConsumable) (err error)
	UpdateConsumable(ctx context.Context, cmd cports.UpdateConsumable) (err error)
	DeleteConsumable(ctx context.Context, uid string) (err error)
	GetConsumableByUID(ctx context.Context, uid string) (result cports.GetConsumable, err error)
	ListConsumable(ctx context.Context, filter cports.Filter) (result []cports.GetConsumable, total int64, err error)
}
