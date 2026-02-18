package ports

import (
	"context"

	"github.com/Logstay/tech-challange-execution-service/pkg/null"
)

type ServiceDefinition interface {
	SaveServiceDefinition(ctx context.Context, cmd SaveServiceDefinition) (err error)
	UpdateServiceDefinition(ctx context.Context, cmd UpdateServiceDefinition) (err error)
	GetServiceDefinitionByUID(ctx context.Context, uid string) (result GetServiceDefinition, err error)
	ListServiceDefinition(ctx context.Context, filter Filter) (result []GetServiceDefinition, err error)
	DeleteServiceDefinition(ctx context.Context, uid string) (err error)
}

type SaveServiceDefinition struct {
	Name        null.String
	Description null.String
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}

type GetServiceDefinition struct {
	ID          null.Int
	UID         null.String
	Name        null.String
	Description null.String
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}

type UpdateServiceDefinition struct {
	UID         null.String
	Name        null.String
	Description null.String
	Price       null.Float
	CreatedAt   null.DateTime
	UpdatedAt   null.DateTime
}
