package internal

import (
	"context"

	"github.com/Logstay/tech-challange-execution-service/config"
	"github.com/Logstay/tech-challange-execution-service/internal/adapters/kafka"
	consumabledao "github.com/Logstay/tech-challange-execution-service/internal/adapters/postgres/consumable"
	inventorydao "github.com/Logstay/tech-challange-execution-service/internal/adapters/postgres/inventory"
	servicedefinitiondao "github.com/Logstay/tech-challange-execution-service/internal/adapters/postgres/servicedefinition"
	consumableservice "github.com/Logstay/tech-challange-execution-service/internal/application/consumable"
	servicedefinitionservice "github.com/Logstay/tech-challange-execution-service/internal/application/servicedefinition"
	eventprocess "github.com/Logstay/tech-challange-execution-service/internal/domain/inventory/events/process"
	cPorts "github.com/Logstay/tech-challange-execution-service/internal/ports/cache"
	database "github.com/Logstay/tech-challange-execution-service/pkg/postgres"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
)

type (
	ConsumableSvc = consumableservice.Service
	ServiceDefSvc = servicedefinitionservice.Service
)

type ServiceBuilder struct {
	dbFactory database.ConnectionFactory
	tracer    *tracing.Provider
	cfg       *config.Config
	cache     cPorts.Cache
}

type ServiceFactory struct {
	*ConsumableSvc
	*ServiceDefSvc
}

func NewServiceBuilder(dbFactory database.ConnectionFactory, tracer *tracing.Provider, cfg *config.Config, cache cPorts.Cache) *ServiceBuilder {
	return &ServiceBuilder{
		dbFactory: dbFactory,
		tracer:    tracer,
		cfg:       cfg,
		cache:     cache,
	}
}

func (b *ServiceBuilder) Build() (*ServiceFactory, error) {
	db := b.dbFactory.NewConnection()

	consumableDAO := consumabledao.NewDAO(db, b.tracer)
	inventoryDAO := inventorydao.NewDAO(db, b.tracer)
	serviceDefDAO := servicedefinitiondao.NewDAO(db, b.tracer)

	consumableRepo := consumabledao.NewRepository(consumableDAO, inventoryDAO, b.tracer)
	serviceDefRepo := servicedefinitiondao.NewRepository(serviceDefDAO, b.tracer)

	consumableSvc := consumableservice.ConsumableService(consumableservice.Dependencies{
		Repository: consumableRepo,
		Tracer:     b.tracer.NewTracer("consumable_service"),
	})

	serviceDefSvc := servicedefinitionservice.ServiceDefinitionService(servicedefinitionservice.Dependencies{
		Repository: serviceDefRepo,
		Tracer:     b.tracer.NewTracer("servicedefinition_service"),
	})

	return &ServiceFactory{
		ConsumableSvc: consumableSvc,
		ServiceDefSvc: serviceDefSvc,
	}, nil
}

func (b *ServiceBuilder) StartEvent(ctx context.Context) error {
	broker, err := kafka.New(kafka.Config{
		ProducerHosts: b.cfg.Kafka.ProducerHosts,
		ConsumerHosts: b.cfg.Kafka.ConsumerHosts,
		Username:      b.cfg.Kafka.Username,
		Password:      b.cfg.Kafka.Password,
	}, b.cfg.Logger)
	if err != nil {
		return err
	}

	processBillingEvent, err := eventprocess.New(broker, b.tracer, b.cfg.Logger)
	if err != nil {
		return err
	}

	go func() {
		if err := processBillingEvent.Handle(ctx); err != nil {
			b.cfg.Logger.Error("Error starting billing event processor", err)
		}
	}()

	select {}
}
