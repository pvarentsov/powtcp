package main

import (
	"os"

	"github.com/pvarentsov/powtcp/internal/app/client"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/config"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/log"
	"github.com/pvarentsov/powtcp/internal/pkg/service"
)

func main() {
	op := "server.main"

	logger := log.New(log.Opts{
		Level: log.LevelDebug,
		Json:  false,
	})

	config, err := config.ParseByFlag("config")
	if err != nil {
		logger.Error(err.Error(), "op", op)
		os.Exit(1)
	}

	service := service.NewClient(service.ClientOpts{
		Config: newConfigService(config),
		Logger: logger,
	})

	err = client.Connect(client.Opts{
		Config:  newConfigClient(config),
		Logger:  logger,
		Service: service,
	})
	if err != nil {
		logger.Error(err.Error(), "op", op)
		os.Exit(1)
	}
}
