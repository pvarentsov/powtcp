package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pvarentsov/powtcp/internal/app/server"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/cache"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/config"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/log"
	"github.com/pvarentsov/powtcp/internal/pkg/lib/tcp"
	"github.com/pvarentsov/powtcp/internal/pkg/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	config, err := config.ParseByFlag("config")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	configService := newConfigService(config)
	configServer := newConfigServer(config)

	logger := log.New(log.Opts{
		Level: log.Level(config.Server.LogLevel),
		Json:  false,
	})

	puzzleCache := cache.New[string, struct{}](ctx, cache.Opts{
		CleanInterval: configService.PuzzleTTL(),
		Logger:        logger,
	})

	service := service.NewServer(service.ServerOpts{
		Config:       configService,
		Logger:       logger,
		PuzzleCache:  puzzleCache,
		ErrorChecker: tcp.NewConnErrorChecker(),
	})

	server, err := server.Listen(ctx, server.Opts{
		Config:  configServer,
		Logger:  logger,
		Service: service,
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	cancel()
	server.Shutdown()
}
