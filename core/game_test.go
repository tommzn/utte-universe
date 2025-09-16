package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type GameSuite struct {
	suite.Suite
	game    *Game
	config  Config
	random  *mockRand
	planets []*Planet
	npcs    []*NPC
	log     *mockLog
}

func TestGameSuite(t *testing.T) {
	suite.Run(t, new(GameSuite))
}

func (s *GameSuite) SetupTest() {
	s.config = Config{TickDuration: 10 * time.Millisecond}
	s.random = &mockRand{seekVal: 0.5, ofVal: 1}
	s.log = &mockLog{}
	s.planets = []*Planet{
		{
			Name:      "TestPlanet",
			Type:      TerraLike,
			Resources: map[ResourceType]int{Iron: 10, Food: 10, Fuel: 10},
			Modifiers: map[ResourceType]float64{},
			Buildings: []*Building{
				{
					Type:       Mine,
					Production: map[ResourceType]int{Iron: 2},
					Modifiers:  map[ResourceType]float64{},
					Level:      1,
				},
				{
					Type:       Farm,
					Production: map[ResourceType]int{Food: 2},
					Modifiers:  map[ResourceType]float64{},
					Level:      1,
				},
				{
					Type:       Refinery,
					Production: map[ResourceType]int{Fuel: 2},
					Modifiers:  map[ResourceType]float64{},
					Level:      1,
				},
			},
		},
		{
			Name:      "SecondPlanet",
			Type:      Desert,
			Resources: map[ResourceType]int{Iron: 5, Food: 5, Fuel: 5},
			Modifiers: map[ResourceType]float64{},
			Buildings: []*Building{
				{
					Type:       Mine,
					Production: map[ResourceType]int{Iron: 1},
					Modifiers:  map[ResourceType]float64{},
					Level:      1,
				},
			},
		},
	}
	s.npcs = []*NPC{
		{
			Name:    "TestNPC",
			Credits: 100,
			Cargo:   map[ResourceType]int{Iron: 0, Food: 0, Fuel: 0},
			Offer:   map[ResourceType]int{Iron: 1, Food: 1, Fuel: 1},
		},
	}
	s.game = NewGameService(s.config, s.random, s.log, s.planets, s.npcs)
	s.game.log = s.log
}

func (s *GameSuite) TestGameLoopCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		s.game.GameLoop(ctx)
		close(done)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case <-done:
		// GameLoop exited as expected
	case <-time.After(100 * time.Millisecond):
		s.Fail("GameLoop did not exit on cancel")
	}
}

func (s *GameSuite) TestGameLoopResourceProduction() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		s.game.GameLoop(ctx)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	// After loop, resources should be produced at least once for both planets
	for _, planet := range s.planets {
		s.GreaterOrEqual(planet.Resources[Iron], 4)
		s.GreaterOrEqual(planet.Resources[Food], 4)
		s.GreaterOrEqual(planet.Resources[Fuel], 4)
	}
}
