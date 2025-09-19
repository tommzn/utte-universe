package core

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	pb "github.com/tommzn/utte-universe/core/proto"
	"google.golang.org/grpc/metadata"
)

type UniverseServerTestSuite struct {
	suite.Suite
	log *mockLog
}

func TestUniverseServerTestSuite(t *testing.T) {
	suite.Run(t, new(UniverseServerTestSuite))
}

func (suite *UniverseServerTestSuite) SetupTest() {
	suite.log = &mockLog{}
}

func (suite *UniverseServerTestSuite) TestGetPlanets() {
	game := &Game{
		Planets: []*Planet{
			{
				Name: "Earth",
				Type: PlanetType(1),
				Resources: map[ResourceType]int{
					ResourceType(1): 100,
				},
				Modifiers: map[ResourceType]float64{
					ResourceType(1): 1.5,
				},
				Buildings: []*Building{
					{
						Type:       BuildingType(1),
						Level:      2,
						Production: map[ResourceType]int{ResourceType(1): 10},
						Modifiers:  map[ResourceType]float64{ResourceType(1): 2.0},
						BuildCost:  map[ResourceType]int{ResourceType(1): 50},
					},
				},
				Owner: &NPC{Name: "NPC1"},
			},
		},
	}
	server := &UniverseServer{Game: game, Log: suite.log}
	resp, err := server.GetPlanets(context.Background(), &pb.Empty{})
	suite.NoError(err)
	suite.Len(resp.Planets, 1)
	suite.Equal("Earth", resp.Planets[0].Name)
}

func (suite *UniverseServerTestSuite) TestGetNPCs() {
	game := &Game{
		NPCs: []*NPC{
			{
				Name:                 "NPC1",
				Offer:                map[ResourceType]int{ResourceType(1): 5},
				Credits:              100,
				Cargo:                map[ResourceType]int{ResourceType(1): 20},
				MaxCargo:             50,
				ColonizationCooldown: time.Now(),
			},
		},
	}
	server := &UniverseServer{Game: game, Log: suite.log}
	resp, err := server.GetNPCs(context.Background(), &pb.Empty{})
	suite.NoError(err)
	suite.Len(resp.Npcs, 1)
	suite.Equal("NPC1", resp.Npcs[0].Name)
}

func (suite *UniverseServerTestSuite) TestPlanetToProto() {
	planet := &Planet{
		Name: "Mars",
		Type: PlanetType(2),
		Resources: map[ResourceType]int{
			ResourceType(2): 200,
		},
		Modifiers: map[ResourceType]float64{
			ResourceType(2): 2.5,
		},
		Buildings: []*Building{
			{
				Type:       BuildingType(2),
				Level:      3,
				Production: map[ResourceType]int{ResourceType(2): 20},
				Modifiers:  map[ResourceType]float64{ResourceType(2): 3.0},
				BuildCost:  map[ResourceType]int{ResourceType(2): 100},
			},
		},
		Owner: &NPC{Name: "NPC2"},
	}
	proto := planetToProto(planet)
	suite.Equal("Mars", proto.Name)
	suite.NotEmpty(proto.Type)
	suite.Len(proto.Buildings, 1)
	suite.NotNil(proto.Owner)
	suite.Equal("NPC2", proto.Owner.Name)
}

func (suite *UniverseServerTestSuite) TestNPCToProto() {
	npc := &NPC{
		Name:                 "NPC3",
		Offer:                map[ResourceType]int{ResourceType(3): 15},
		Credits:              300,
		Cargo:                map[ResourceType]int{ResourceType(3): 30},
		MaxCargo:             60,
		ColonizationCooldown: time.Now(),
	}
	proto := npcToProto(npc)
	suite.Equal("NPC3", proto.Name)
	suite.Equal(int32(300), proto.Credits)
	suite.Equal(int32(60), proto.MaxCargo)
}

func (suite *UniverseServerTestSuite) TestEventToProto() {
	event := &Event{
		Name:           "Boost",
		Target:         1,
		TargetPlanet:   &Planet{Name: "Venus"},
		TargetBuilding: &Building{Type: BuildingType(4)},
		ResourceBoost:  map[ResourceType]float64{ResourceType(4): 4.4},
		Duration:       10,
		RemainingTicks: 5,
	}
	proto := eventToProto(event)
	suite.Equal("Boost", proto.Name)
	suite.Equal("Venus", proto.TargetPlanet)
	suite.NotEmpty(proto.TargetBuilding)
	suite.Equal(int32(10), proto.Duration)
	suite.Equal(int32(5), proto.RemainingTicks)
}

func (suite *UniverseServerTestSuite) TestPlanetToProtoNilOwner() {
	planet := &Planet{
		Name: "NoOwner",
		Type: PlanetType(1),
	}
	proto := planetToProto(planet)
	suite.Nil(proto.Owner)
}

func (suite *UniverseServerTestSuite) TestNPCToProtoZeroCooldown() {
	npc := &NPC{
		Name: "NoCooldown",
	}
	proto := npcToProto(npc)
	suite.Empty(proto.ColonizationCooldown)
}

func (suite *UniverseServerTestSuite) TestEventToProtoNilTargets() {
	event := &Event{
		Name:          "NoTargets",
		ResourceBoost: map[ResourceType]float64{},
	}
	proto := eventToProto(event)
	suite.Empty(proto.TargetPlanet)
	suite.Empty(proto.TargetBuilding)
}

func (suite *UniverseServerTestSuite) TestStartGRPCServerError() {
	game := &Game{}
	log := &mockLog{}
	// Use an invalid address to force error
	_, _, err := NewGRPCServer(game, "invalid_addr", log)
	suite.Error(err)
}

type mockStream struct {
	recvCalls  int
	sendCalls  int
	recvCmds   []*pb.ClientCommand
	sentStates []*pb.UniverseState
	closed     bool
	ctx        context.Context
}

func (m *mockStream) Recv() (*pb.ClientCommand, error) {
	if m.recvCalls < len(m.recvCmds) {
		cmd := m.recvCmds[m.recvCalls]
		m.recvCalls++
		return cmd, nil
	}
	m.closed = true
	return nil, errors.New("stream closed")
}

func (m *mockStream) Send(state *pb.UniverseState) error {
	m.sentStates = append(m.sentStates, state)
	m.sendCalls++
	return nil
}

func (m *mockStream) Context() context.Context {
	if m.ctx == nil {
		return context.Background()
	}
	return m.ctx
}

func (m *mockStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockStream) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockStream) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockStream) SendHeader(md metadata.MD) error {
	return nil
}

func (m *mockStream) SetTrailer(md metadata.MD) {
	// no-op
}

// Implement BidiStreamingServer interface marker method
func (m *mockStream) BidiStreamingServer() {}

type mockStreamErrorSend struct {
	recvCalls int
	recvCmds  []*pb.ClientCommand
	ctx       context.Context
}

func (m *mockStreamErrorSend) SetHeader(md metadata.MD) error {
	return nil
}

func (m *mockStreamErrorSend) SetTrailer(md metadata.MD) {
	// no-op
}

func (m *mockStreamErrorSend) Recv() (*pb.ClientCommand, error) {
	if m.recvCalls < len(m.recvCmds) {
		cmd := m.recvCmds[m.recvCalls]
		m.recvCalls++
		return cmd, nil
	}
	return nil, errors.New("stream closed")
}

func (m *mockStreamErrorSend) Send(state *pb.UniverseState) error {
	return errors.New("send error")
}

func (m *mockStreamErrorSend) Context() context.Context {
	if m.ctx == nil {
		return context.Background()
	}
	return m.ctx
}

func (m *mockStreamErrorSend) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockStreamErrorSend) RecvMsg(msg interface{}) error {
	return nil
}

func (m *mockStreamErrorSend) SendHeader(md metadata.MD) error {
	return nil
}

// Implement BidiStreamingServer interface marker method
func (m *mockStreamErrorSend) BidiStreamingServer() {}

func (suite *UniverseServerTestSuite) TestStreamUniverseStateSubscribePauseResumeUnsubscribe() {
	// Setup channels
	planetCh := make(chan []*Planet, 1)
	npcCh := make(chan []*NPC, 1)
	eventCh := make(chan []*Event, 1)

	// Prepare updates
	planetCh <- []*Planet{{Name: "Earth"}}
	npcCh <- []*NPC{{Name: "NPC1"}}
	eventCh <- []*Event{{Name: "Event1"}}

	game := &Game{
		planetUpdates: planetCh,
		npcUpdates:    npcCh,
		eventUpdates:  eventCh,
	}

	server := &UniverseServer{Game: game, Log: suite.log}

	// Prepare commands: SUBSCRIBE, PAUSE, RESUME, UNSUBSCRIBE, then close
	cmds := []*pb.ClientCommand{
		{Type: pb.ClientCommand_SUBSCRIBE},
		{Type: pb.ClientCommand_PAUSE},
		{Type: pb.ClientCommand_RESUME},
		{Type: pb.ClientCommand_UNSUBSCRIBE},
	}
	stream := &mockStream{recvCmds: cmds}

	// Run the method in a goroutine to avoid blocking
	done := make(chan error, 1)
	go func() {
		err := server.StreamUniverseState(stream)
		// Only treat error as failure if it's not "stream closed"
		if err != nil && err.Error() != "stream closed" {
			done <- err
		} else {
			done <- nil
		}
	}()

	// Wait for goroutine to finish
	err := <-done
	suite.True(stream.closed)
	suite.NoError(err)
	// Only one update should be sent (on SUBSCRIBE and RESUME, but only one set of updates in channel)
	suite.GreaterOrEqual(stream.sendCalls, 1)
	suite.Equal("Earth", stream.sentStates[0].Planets.Planets[0].Name)
	suite.Equal("NPC1", stream.sentStates[0].Npcs.Npcs[0].Name)
	suite.Equal("Event1", stream.sentStates[0].Events[0].Name)
}

func (suite *UniverseServerTestSuite) TestStreamUniverseStateErrorOnRecv() {
	game := &Game{}
	server := &UniverseServer{Game: game, Log: suite.log}
	stream := &mockStream{recvCmds: []*pb.ClientCommand{}}
	err := server.StreamUniverseState(stream)
	suite.Error(err)
}

func (suite *UniverseServerTestSuite) TestStreamUniverseStateErrorOnSend() {
	planetCh := make(chan []*Planet, 1)
	npcCh := make(chan []*NPC, 1)
	eventCh := make(chan []*Event, 1)
	planetCh <- []*Planet{{Name: "Earth"}}
	npcCh <- []*NPC{{Name: "NPC1"}}
	eventCh <- []*Event{{Name: "Event1"}}
	game := &Game{
		planetUpdates: planetCh,
		npcUpdates:    npcCh,
		eventUpdates:  eventCh,
	}
	server := &UniverseServer{Game: game, Log: suite.log}
	// Stream that returns error on Send
	stream := &mockStreamErrorSend{
		recvCmds: []*pb.ClientCommand{{Type: pb.ClientCommand_SUBSCRIBE}},
	}
	err := server.StreamUniverseState(stream)
	suite.Error(err)
}
