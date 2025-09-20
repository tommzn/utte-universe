package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tommzn/go-log"
	pb "github.com/tommzn/utte-universe/core/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type UIBBackend struct {
	gameClient pb.UniverseServiceClient
	logger     log.Logger
}

func NewUIBackend(addr string, logger log.Logger) (*UIBBackend, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial game backend: %w", err)
	}
	client := pb.NewUniverseServiceClient(conn)
	return &UIBBackend{gameClient: client, logger: logger}, nil
}

func (u *UIBBackend) flushLogsPeriodically(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			u.logger.Flush()
		}
	}
}

func (u *UIBBackend) handleWebsocket(w http.ResponseWriter, r *http.Request) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		u.logger.Errorf("WebSocket upgrade failed: %v", err)
		http.Error(w, "failed to upgrade", http.StatusBadRequest)
		return
	}
	defer ws.Close()
	u.logger.Info("WebSocket connection established")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go u.flushLogsPeriodically(ctx)

	stream, err := u.gameClient.StreamUniverseState(ctx)
	if err != nil {
		u.logger.Errorf("Failed to open gRPC stream: %v", err)
		return
	}
	u.logger.Info("gRPC stream to backend opened")

	// Forward frontend → backend
	go func() {
		for {
			var msg map[string]string
			if err := ws.ReadJSON(&msg); err != nil {
				u.logger.Errorf("WebSocket read error: %v", err)
				cancel()
				return
			}
			if command, ok := msg["command"]; ok {
				u.logger.Debugf("Received command from frontend: %s", command)
				if err := stream.Send(&pb.ClientCommand{Type: commandTypeFromString(command)}); err != nil {
					u.logger.Errorf("Failed to send command to backend: %v", err)
					cancel()
					return
				}
			}
		}
	}()

	// Forward backend → frontend
	for {
		update, err := stream.Recv()
		if err != nil {
			u.logger.Errorf("gRPC stream recv error: %v", err)
			return
		}
		u.logger.Debug("Sending update to frontend")
		if err := ws.WriteJSON(update); err != nil {
			u.logger.Errorf("WebSocket write error: %v, content: %+v", err, update)
			return
		}
	}
}

func commandTypeFromString(cmd string) pb.ClientCommand_CommandType {
	switch cmd {
	case "SUBSCRIBE":
		return pb.ClientCommand_SUBSCRIBE
	case "PAUSE":
		return pb.ClientCommand_PAUSE
	case "RESUME":
		return pb.ClientCommand_RESUME
	case "UNSUBSCRIBE":
		return pb.ClientCommand_UNSUBSCRIBE
	default:
		return pb.ClientCommand_SUBSCRIBE // or a default/fallback value
	}
}
