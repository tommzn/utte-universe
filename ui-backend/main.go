package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger := bootstrap()
	defer logger.Flush()

	backendAddr := os.Getenv("GAME_BACKEND_ADDR")
	if backendAddr == "" {
		backendAddr = "localhost:8081"
	}

	ui, err := NewUIBackend(backendAddr, logger)
	if err != nil {
		logger.Errorf("Failed to connect to game backend: %v", err)
		os.Exit(1)
	}
	ui.logger = logger // inject logger from bootstrap

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ui.handleWebsocket)

	httpPort := os.Getenv("HTTP_PORT")
	healthServer := NewHealthServer(":"+httpPort, logger)
	go healthServer.Start()

	srv := &http.Server{Addr: ":" + httpPort, Handler: mux}
	go func() {
		logger.Infof("UI Backend listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("HTTP server error: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down UI backend...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server forced to shutdown: %v", err)
		os.Exit(1)
	}

	healthServer.Shutdown(5 * time.Second)

	logger.Info("UI backend exited cleanly")
}
