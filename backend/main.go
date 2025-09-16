package main

import (
	"context"
	"fmt"
	"net/http"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// Metrics (simple example, extend with Prometheus later)
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "active_events %d\n", len(game.ActiveEvents))
		fmt.Fprintf(w, "planets %d\n", len(game.Planets))
		fmt.Fprintf(w, "npcs %d\n", len(game.NPCs))
	})

	// Example endpoint: list planets
	mux.HandleFunc("/planets", func(w http.ResponseWriter, r *http.Request) {
		for _, p := range game.Planets {
			owner := "none"
			if p.Owner != nil {
				owner = p.Owner.Name
			}
			fmt.Fprintf(w, "Planet: %s (Type: %s, Owner: %s, Resources: %+v)\n",
				p.Name, p.Type.String(), owner, p.Resources)
		}
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		logger.Infof("Backend listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("ListenAndServe error: %v", err)
		}
	}()

	// Handle OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down...")

	// Cancel context → stops GameLoop
	cancel()

	// Graceful server shutdown
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Exited cleanly")
}

func AsGameLogger(logger log.Logger) core.Log {
	return core.NewCustomLogger(logger)
}
