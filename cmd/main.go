package main

import (
	"context"
	"log"
	"os"

	"github.com/Logstay/tech-challange-execution-service/cmd/initialize"
)

func main() {
	ctx := context.Background()

	app, err := initialize.Setup(ctx)
	if err != nil {
		log.Fatalf("Falha no bootstrap da aplicação: %v", err)
	}

	defer func() {
		if shutdownErr := app.Shutdown(ctx); shutdownErr != nil {
			app.Config.Logger.Error("Erro durante o shutdown", shutdownErr)
		}
		os.Exit(0)
	}()

	serverManager, err := initialize.NewServerManager(ctx, app)
	if err != nil {
		app.Config.Logger.Error("Falha ao criar server manager", err)
		return
	}

	serverManager.StartEventProcessors(ctx)

	if err := serverManager.Start(); err != nil {
		app.Config.Logger.Error("Erro na execução do servidor", err)
		return
	}
}
