package ports

import (
	"context"

	modelservicedefinition "github.com/Logstay/tech-challange-wherehouse-service/internal/domain/servicedefinition"
)

type Repository interface {
	Save(ctx context.Context, in modelservicedefinition.ServiceDefinition) (err error)
	Get(ctx context.Context, uid string) (result GetServiceDefinition, err error)
	List(ctx context.Context, filter Filter) (result []GetServiceDefinition, total int64, err error)
	Update(ctx context.Context, uid string, in modelservicedefinition.ServiceDefinition) (err error)
	Delete(ctx context.Context, uid string) (err error)
	Exist(ctx context.Context, name string) (result bool, err error)
}
