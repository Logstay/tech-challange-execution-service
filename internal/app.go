package services

import (
	"context"

	"github.com/Logstay/tech-challange-wherehouse-service/config"
	// consumabledao "github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/postgres/consumable"
	billingadapter "github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/billing/mercado_pago"
	"github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/kafka"
	billingservice "github.com/Logstay/tech-challange-wherehouse-service/internal/application/billing"
	eventprocessbilling "github.com/Logstay/tech-challange-wherehouse-service/internal/domain/billing/events/process-billing"
	bPorts "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/billing"
	cPorts "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/cache"
	messagebroker "github.com/Logstay/tech-challange-wherehouse-service/internal/ports/message-broker"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/tracing"
)

type ServiceFactory struct {
	tracer *tracing.Provider
	cfg    *config.Config
	cache  cPorts.Cache
}

func NewServiceFactory(
	tracer *tracing.Provider,
	cfg *config.Config,
	cache cPorts.Cache,
) *ServiceFactory {
	return &ServiceFactory{
		tracer: tracer,
		cfg:    cfg,
		cache:  cache,
	}
}

func (f *ServiceFactory) BuildApp(ctx context.Context) (service *billingservice.Service, err error) {

	billingPayment, err := f.buildBillingAdapter()
	if err != nil {
		f.cfg.Logger.Error("Error building billing adapter", err)
		return
	}

	service = billingservice.NewService(billingservice.Dependencies{
		Tracer:         f.tracer.NewTracer("billing_service"),
		BillingAdapter: billingPayment,
	})

	return

}

func buildDatavizMessageBroker(cfg *config.Config) (broker messagebroker.MessageBroker, err error) {
	return kafka.New(kafka.Config{
		ProducerHosts: cfg.Kafka.ProducerHosts,
		ConsumerHosts: cfg.Kafka.ConsumerHosts,
		Username:      cfg.Kafka.Username,
		Password:      cfg.Kafka.Password,
	}, cfg.Logger)
}

func (f *ServiceFactory) StartEvent(ctx context.Context) (err error) {
	broker, err := buildDatavizMessageBroker(f.cfg)
	if err != nil {
		return
	}

	processBillingEvent, err := eventprocessbilling.New(broker, f.tracer, f.cfg.Logger)
	if err != nil {
		return
	}

	go func() {
		err = processBillingEvent.Handle(ctx)
		if err != nil {
			f.cfg.Logger.Error("Error starting billing event processor", err)
		}
	}()

	select {}

}

func (f *ServiceFactory) buildBillingAdapter() (adapter bPorts.MercadoPagoBilling, err error) {

	return billingadapter.NewMercadoPagoPayment(f.cfg.MercadoPago.AccessToken, f.cfg.Logger)
}
