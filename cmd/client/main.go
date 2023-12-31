package main

import (
	"fmt"
	"os"

	"github.com/pvarentsov/powtcp/internal/app/client"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/config"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/log"
	"github.com/pvarentsov/powtcp/internal/pkg/service"
)

func main() {
	config, err := config.Parse("config")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	configService := newConfigService(config)
	configClient := newConfigClient(config)

	logger := log.New(log.Opts{
		Level: log.Level(config.Client.LogLevel),
		Json:  config.Client.LogJson,
	})

	service := service.NewClient(service.ClientOpts{
		Config: configService,
		Logger: logger,
	})

	logger.Debug("client configured",
		"server_address", configClient.ServerAddress(),
		"puzzle_compute_max_attempts", configService.PuzzleComputeMaxAttempts(),
	)

	err = client.Connect(client.Opts{
		Config:  configClient,
		Logger:  logger,
		Service: service,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
