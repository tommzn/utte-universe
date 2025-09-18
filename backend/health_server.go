package main

import (
	"context"
	"net/http"
	"time"

	"github.com/tommzn/go-log"
)

type HealthServer struct {
	server *http.Server
	logger log.Logger
	done   chan struct{}
}

func NewHealthServer(addr string, logger log.Logger) *HealthServer {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	return &HealthServer{
		server: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
		logger: logger,
		done:   make(chan struct{}),
	}
}

func (hs *HealthServer) Start() {
	hs.logger.Infof("Health server listening on %s", hs.server.Addr)
	if err := hs.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		hs.logger.Errorf("Health server error: %v", err)
	}
	close(hs.done)
}

func (hs *HealthServer) Shutdown(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := hs.server.Shutdown(ctx); err != nil {
		hs.logger.Errorf("Health server forced to shutdown: %v", err)
	}
	<-hs.done
}
