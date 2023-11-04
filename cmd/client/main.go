package main

import (
	"os"

	"github.com/pvarentsov/powtcp/internal/app/client"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/log"
	"github.com/pvarentsov/powtcp/internal/pkg/service"
)

func main() {
	op := "server.main"

	logger := log.New(log.Opts{
		Level: log.LevelDebug,
		Json:  false,
	})

	service := service.NewClient(service.ClientOpts{
		Logger: logger,
	})

	err := client.Connect(client.Opts{
		Address: ":8080",
		Logger:  logger,
		Service: service,
	})
	if err != nil {
		logger.Error(err.Error(), "op", op)
		os.Exit(1)
	}
}
