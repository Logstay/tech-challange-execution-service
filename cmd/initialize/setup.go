package initialize

import (
	"context"
	"errors"
	"fmt"

	"github.com/Logstay/tech-challange-execution-service/config"
	app "github.com/Logstay/tech-challange-execution-service/internal"
	interfaces "github.com/Logstay/tech-challange-execution-service/internal/interfaces/services"
	cPorts "github.com/Logstay/tech-challange-execution-service/internal/ports/cache"
	"github.com/Logstay/tech-challange-execution-service/pkg/cache"
	"github.com/Logstay/tech-challange-execution-service/pkg/constant"
	database "github.com/Logstay/tech-challange-execution-service/pkg/postgres"
	"github.com/Logstay/tech-challange-execution-service/pkg/tracing"
	_ "github.com/lib/pq"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Application struct {
	Config          *config.Config
	TracingProvider *tracing.Provider
	Cache           cPorts.Cache
	DB              database.ConnectionFactory
	ServiceFactory  interfaces.ServiceFactory
	NewRelicApp     *newrelic.Application
	Context         context.Context
}

func Setup(ctx context.Context) (*Application, error) {
	appInstance := &Application{
		Context: ctx,
	}

	if err := appInstance.loadConfig(); err != nil {
		return nil, fmt.Errorf("falha ao carregar configurações: %w", err)
	}

	if err := appInstance.initNewRelic(); err != nil {
		appInstance.Config.Logger.Error("Erro ao inicializar New Relic (continuando sem ele)", err)
	}

	if err := appInstance.initTracing(); err != nil {
		appInstance.Config.Logger.Error("Erro ao inicializar tracing provider (continuando sem ele)", err)
	}

	if err := appInstance.initDatabase(); err != nil {
		return nil, fmt.Errorf("falha ao inicializar banco de dados: %w", err)
	}

	appInstance.initCache()

	appInstance.prepareContext()

	appInstance.initServiceFactory()

	appInstance.Config.Logger.Info("Bootstrap concluído com sucesso")

	return appInstance, nil
}

func (a *Application) Shutdown(ctx context.Context) error {
	a.Config.Logger.Info("Iniciando shutdown da aplicação...")

	if a.TracingProvider != nil {
		a.TracingProvider.Close(ctx)
	}

	if a.Cache != nil {
		a.Cache.Close()
	}

	if a.DB != nil {
		if err := a.DB.Close(); err != nil {
			a.Config.Logger.Error("Failed to close database connection", err)
		}
	}

	a.Config.Logger.Info("Shutdown concluído")
	return nil
}

func (a *Application) loadConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	a.Config = cfg
	return nil
}

func (a *Application) initNewRelic() error {
	if a.Config.NewRelic.LicenseKey == "" {
		return errors.New("new Relic license key não configurada")
	}

	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(a.Config.NewRelic.AppName),
		newrelic.ConfigLicense(a.Config.NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(true),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		return err
	}

	a.NewRelicApp = nrApp
	return nil
}

func (a *Application) initTracing() error {
	tracingProvider, err := tracing.NewTracingProvider(a.Config)
	if err != nil {
		return err
	}
	a.TracingProvider = tracingProvider
	return nil
}

func (a *Application) initCache() {
	a.Cache = cache.NewMemoryCache()
}

func (a *Application) initDatabase() error {
	db, err := database.NewPostgresConnectionFactory(a.Context, a.Config.Database)
	if err != nil {
		return fmt.Errorf("falha ao conectar ao banco de dados: %w", err)
	}
	a.DB = db
	return nil
}

func (a *Application) prepareContext() {
	a.Context = context.WithValue(a.Context, constant.TransactionName("logger"), a.Config.Logger)
}

func (a *Application) initServiceFactory() {

	factory, err := app.NewServiceBuilder(a.DB, a.TracingProvider, a.Config, a.Cache).Build()
	if err != nil {
		a.Config.Logger.Error("Erro ao criar service factory", err)
		panic(err)
	}
	a.ServiceFactory = factory
}
