package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	app "github.com/Logstay/tech-challange-execution-service/internal"
	handler "github.com/Logstay/tech-challange-execution-service/internal/adapters/http"
	"github.com/Logstay/tech-challange-execution-service/pkg/constant"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
)

type ServerManager struct {
	app        *Application
	httpServer *http.Server
	serverHTTP *gin.Engine
	servers    run.Group
}

func NewServerManager(ctx context.Context, app *Application) (*ServerManager, error) {
	sm := &ServerManager{
		app: app,
	}

	sm.setupHTTPServer()

	sm.setupRunGroups()

	return sm, nil
}

func (sm *ServerManager) setupHTTPServer() {

	sm.serverHTTP = handler.NewServer(
		sm.app.Context,
		sm.app.ServiceFactory,
		sm.app.Config,
		sm.app.TracingProvider,
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
	factory := app.NewServiceBuilder(sm.app.DB, sm.app.TracingProvider, sm.app.Config, sm.app.Cache)

	sm.app.Config.Logger.Info("Iniciando processadores de eventos...")
	go func() {
		if err := factory.StartEvent(ctx); err != nil {
			sm.app.Config.Logger.Error("Erro no processador de eventos de billing", err)
		}
	}()
}
