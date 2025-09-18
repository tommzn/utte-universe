package core

import (
	"context"
	"net"
	"time"

	pb "github.com/tommzn/utte-universe/core/proto"
	"google.golang.org/grpc"
)

type UniverseServer struct {
	pb.UnimplementedUniverseServiceServer
	Game *Game
	Log  Log
}

func (s *UniverseServer) GetPlanets(ctx context.Context, in *pb.Empty) (*pb.PlanetList, error) {
	s.Log.Info("Received GetPlanets request")
	planets := make([]*pb.Planet, 0, len(s.Game.Planets))
	for _, p := range s.Game.Planets {
		planets = append(planets, planetToProto(p))
	}
	s.Log.Debug("Returning %d planets", len(planets))
	return &pb.PlanetList{Planets: planets}, nil
}

func (s *UniverseServer) GetNPCs(ctx context.Context, in *pb.Empty) (*pb.NPCList, error) {
	s.Log.Info("Received GetNPCs request")
	npcs := make([]*pb.NPC, 0, len(s.Game.NPCs))
	for _, n := range s.Game.NPCs {
		npcs = append(npcs, npcToProto(n))
	}
	s.Log.Debug("Returning %d NPCs", len(npcs))
	return &pb.NPCList{Npcs: npcs}, nil
}

func (s *UniverseServer) StreamUniverseState(stream pb.UniverseService_StreamUniverseStateServer) error {
	s.Log.Info("Started StreamUniverseState")
	subscribed := false
	paused := false

	for {
		select {
		case <-stream.Context().Done():
			s.Log.Info("StreamUniverseState context cancelled")
			return stream.Context().Err()
		default:
			cmd, err := stream.Recv()
			if err != nil {
				s.Log.Error("StreamUniverseState closed or errored: %v", err)
				return err
			}

			s.Log.Debug("Received client command: %v", cmd.Type)
			switch cmd.Type {
			case pb.ClientCommand_SUBSCRIBE:
				subscribed = true
				paused = false
				s.Log.Info("Client subscribed to universe state stream")
			case pb.ClientCommand_PAUSE:
				paused = true
				s.Log.Info("Client paused universe state stream")
			case pb.ClientCommand_RESUME:
				paused = false
				s.Log.Info("Client resumed universe state stream")
			case pb.ClientCommand_UNSUBSCRIBE:
				subscribed = false
				s.Log.Info("Client unsubscribed from universe state stream")
			}

			if subscribed && !paused {
			select {
			case planets := <-s.Game.planetUpdates:
				var npcs []*NPC
				var events []*Event
				select {
				case npcs = <-s.Game.npcUpdates:
				default:
					npcs = []*NPC{}
				}
				select {
				case events = <-s.Game.eventUpdates:
				default:
					events = []*Event{}
				}

				s.Log.Debug("Sending universe state update: %d planets, %d NPCs, %d events", len(planets), len(npcs), len(events))

				planetsProto := make([]*pb.Planet, 0, len(planets))
				for _, p := range planets {
					planetsProto = append(planetsProto, planetToProto(p))
				}
				npcsProto := make([]*pb.NPC, 0, len(npcs))
				for _, n := range npcs {
					npcsProto = append(npcsProto, npcToProto(n))
				}
				eventsProto := make([]*pb.Event, 0, len(events))
				for _, e := range events {
					eventsProto = append(eventsProto, eventToProto(e))
				}
				msg := &pb.UniverseState{
					Planets: &pb.PlanetList{Planets: planetsProto},
					Npcs:    &pb.NPCList{Npcs: npcsProto},
					Events:  eventsProto,
				}
				if err := stream.Send(msg); err != nil {
					s.Log.Error("Failed to send universe state: %v", err)
					return err
				}
			default:
				// No planet updates available, continue loop
			}
		}
	}
}

func planetToProto(p *Planet) *pb.Planet {
	resources := make(map[string]int32)
	for k, v := range p.Resources {
		resources[k.String()] = int32(v)
	}
	modifiers := make(map[string]float32)
	for k, v := range p.Modifiers {
		modifiers[k.String()] = float32(v)
	}
	buildings := make([]*pb.Building, 0, len(p.Buildings))
	for _, b := range p.Buildings {
		bRes := make(map[string]int32)
		for k, v := range b.Production {
			bRes[k.String()] = int32(v)
		}
		bMods := make(map[string]float32)
		for k, v := range b.Modifiers {
			bMods[k.String()] = float32(v)
		}
		bCost := make(map[string]int32)
		for k, v := range b.BuildCost {
			bCost[k.String()] = int32(v)
		}
		buildings = append(buildings, &pb.Building{
			Type:       b.Type.String(),
			Level:      int32(b.Level),
			Production: bRes,
			Modifiers:  bMods,
			BuildCost:  bCost,
		})
	}
	var owner *pb.NPC
	if p.Owner != nil {
		owner = npcToProto(p.Owner)
	}
	return &pb.Planet{
		Name:      p.Name,
		Type:      p.Type.String(),
		Resources: resources,
		Modifiers: modifiers,
		Buildings: buildings,
		Owner:     owner,
	}
}

func npcToProto(n *NPC) *pb.NPC {
	offer := make(map[string]int32)
	for k, v := range n.Offer {
		offer[k.String()] = int32(v)
	}
	cargo := make(map[string]int32)
	for k, v := range n.Cargo {
		cargo[k.String()] = int32(v)
	}
	var cooldown string
	if !n.ColonizationCooldown.IsZero() {
		cooldown = n.ColonizationCooldown.Format(time.RFC3339)
	}
	return &pb.NPC{
		Name:                 n.Name,
		Offer:                offer,
		Credits:              int32(n.Credits),
		Cargo:                cargo,
		MaxCargo:             int32(n.MaxCargo),
		ColonizationCooldown: cooldown,
	}
}

func eventToProto(e *Event) *pb.Event {
	resourceBoost := make(map[string]float32)
	for k, v := range e.ResourceBoost {
		resourceBoost[k.String()] = float32(v)
	}
	var targetPlanet, targetBuilding string
	if e.TargetPlanet != nil {
		targetPlanet = e.TargetPlanet.Name
	}
	if e.TargetBuilding != nil {
		targetBuilding = e.TargetBuilding.Type.String()
	}
	return &pb.Event{
		Name:           e.Name,
		Target:         int32(e.Target),
		TargetPlanet:   targetPlanet,
		TargetBuilding: targetBuilding,
		ResourceBoost:  resourceBoost,
		Duration:       int32(e.Duration),
		RemainingTicks: int32(e.RemainingTicks),
	}
}

func StartGRPCServer(game *Game, addr string, log Log) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to listen on %s: %v", addr, err)
		return err
	}
	grpcServer := grpc.NewServer()
	pb.RegisterUniverseServiceServer(grpcServer, &UniverseServer{Game: game, Log: log})
	log.Info("gRPC server started on %s", addr)
	return grpcServer.Serve(lis)
}
