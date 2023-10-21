package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/pvarentsov/powtcp/internal/lib/log"
	"github.com/pvarentsov/powtcp/internal/server"
)

func main() {
	op := "main"
	ctx, cancel := context.WithCancel(context.Background())

	logger := log.New(log.Opts{
		Debug: true,
		Json:  false,
	})

	server, err := server.Listen(ctx, server.Opts{
		Address: ":8080",
		Logger:  logger,
	})
	if err != nil {
		logger.Error(err.Error(), "op", op)
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	cancel()
	server.Shutdown()
}
