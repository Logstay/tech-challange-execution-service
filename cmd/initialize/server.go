package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	billingHandler "github.com/Logstay/tech-challange-wherehouse-service/internal/adapters/http"
	billingservice "github.com/Logstay/tech-challange-wherehouse-service/internal/application/billing"
	"github.com/Logstay/tech-challange-wherehouse-service/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
)

type ServerManager struct {
	app            *Application
	billingService *billingservice.Service
	httpServer     *http.Server
	serverHTTP     *gin.Engine
	servers        run.Group
}

func NewServerManager(ctx context.Context, app *Application) (*ServerManager, error) {
	sm := &ServerManager{
		app: app,
	}

	if err := sm.initBillingService(ctx); err != nil {
		return nil, fmt.Errorf("falha ao inicializar billing service: %w", err)
	}

	sm.setupHTTPServer()

	sm.setupRunGroups()

	return sm, nil
}

func (sm *ServerManager) initBillingService(ctx context.Context) error {
	service, err := sm.app.ServiceFactory.BuildApp(ctx)
	if err != nil {
		return err
	}
	sm.billingService = service

	return nil
}

func (sm *ServerManager) setupHTTPServer() {

	sm.serverHTTP = billingHandler.NewServer(
		sm.app.Context,
		sm.billingService,
		sm.app.Config,
	).GetRouter()

	sm.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%s", sm.app.Config.Server.Host, sm.app.Config.Server.Port),
		Handler:      sm.serverHTTP,
		WriteTimeout: time.Second * constant.ServerWriteReadTimeout,
		ReadTimeout:  time.Second * constant.ServerWriteReadTimeout,
		IdleTimeout:  time.Second * constant.ServerIdleTimeout,
	}

	sm.app.Config.Logger.Info(fmt.Sprintf(
		"Servidor HTTP configurado em %s:%s",
		sm.app.Config.Server.Host,
		sm.app.Config.Server.Port,
	))
}

func (sm *ServerManager) setupRunGroups() {

	sm.servers.Add(
		func() error {
			sm.app.Config.Logger.Info("Iniciando servidor HTTP...")
			return sm.httpServer.ListenAndServe()
		},
		func(err error) {
			sm.app.Config.Logger.Info("Encerrando servidor HTTP...")
			ctx, cancel := context.WithTimeout(context.Background(), constant.ServerContextTimeout)
			defer cancel()

			if shutdownErr := sm.httpServer.Shutdown(ctx); shutdownErr != nil {
				sm.app.Config.Logger.Error("Erro ao encerrar servidor HTTP", shutdownErr)
			}
		},
	)

	sm.servers.Add(
		func() error {
			return sm.waitForShutdownSignal()
		},
		func(error) {
			sm.app.Config.Logger.Info("Sinal de shutdown recebido")
		},
	)
}

func (sm *ServerManager) waitForShutdownSignal() error {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	sm.app.Config.Logger.Info(fmt.Sprintf("Sinal %v recebido, iniciando shutdown...", sig))

	return nil
}

func (sm *ServerManager) Start() error {
	sm.app.Config.Logger.Info("Servidor iniciado e pronto para receber requisições")
	return sm.servers.Run()
}

func (sm *ServerManager) StartEventProcessors(ctx context.Context) {
	sm.app.Config.Logger.Info("Iniciando processadores de eventos...")
	go func() {
		if err := sm.app.ServiceFactory.StartEvent(ctx); err != nil {
			sm.app.Config.Logger.Error("Erro no processador de eventos de billing", err)
		}
	}()
}
