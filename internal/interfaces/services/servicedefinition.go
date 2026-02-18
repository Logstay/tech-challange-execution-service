package services

import (
	"context"

	sdports "github.com/Logstay/tech-challange-execution-service/internal/ports/servicedefinition"
)

type ServiceDefinitionService interface {
	SaveServiceDefinition(ctx context.Context, cmd sdports.SaveServiceDefinition) (err error)
	UpdateServiceDefinition(ctx context.Context, cmd sdports.UpdateServiceDefinition) (err error)
	DeleteServiceDefinition(ctx context.Context, uid string) (err error)
	GetServiceDefinitionByUID(ctx context.Context, uid string) (result sdports.GetServiceDefinition, err error)
	ListServiceDefinition(ctx context.Context, filter sdports.Filter) (result []sdports.GetServiceDefinition, total int64, err error)
}
