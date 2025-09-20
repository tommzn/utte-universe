package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/tommzn/go-config"
	"github.com/tommzn/go-log"
	"github.com/tommzn/utte-universe/core"
)

func main() {

	conf, _, logger, ctx := bootstrap()
	httpPort := os.Getenv("HTTP_PORT")
	grpcPort := os.Getenv("GRPC_PORT")
	defer logger.Flush()

	gameConfig := &core.Config{}
	if err := gameConfig.LoadFrom(conf); err != nil {
		logger.Error("Failed to load game configuration: %v", err)
		os.Exit(1)
	}

	gameLogger := AsGameLogger(logger)
	rand := core.NewBuiltInRand()

	planet, npcs := core.SeedUniverse(gameConfig.SeedConfig, rand)
	game := core.NewGameService(*gameConfig, rand, gameLogger, planet, npcs)

	gameCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		game.GameLoop(gameCtx)
	}()

	healthServer := NewHealthServer(":"+httpPort, logger)
	go healthServer.Start()

	// Graceful gRPC server setup
	grpcDone := make(chan struct{})
	grpcServer, grpcListener, err := core.NewGRPCServer(game, ":"+grpcPort, gameLogger)
	if err != nil {
		logger.Error("Failed to start gRPC server: %v", err)
		os.Exit(1)
	}
	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			logger.Error("gRPC server error: %v", err)
		}
		close(grpcDone)
	}()

	// Handle OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down...")

	// Cancel context → stops GameLoop
	cancel()

	healthServer.Shutdown(5 * time.Second)

	// Graceful shutdown for gRPC server
	grpcServer.GracefulStop()
	<-grpcDone

	logger.Info("Exited cleanly")
}

func AsGameLogger(logger log.Logger) core.Log {
	return core.NewCustomLogger(logger)
}
