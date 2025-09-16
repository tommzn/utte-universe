package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type NPCSuite struct {
	suite.Suite
	npc     *NPC
	planets []*Planet
	log     *mockLog
}

func TestNPCSuite(t *testing.T) {
	suite.Run(t, new(NPCSuite))
}

func (s *NPCSuite) SetupTest() {
	s.log = &mockLog{}
	s.npc = &NPC{
		Credits:  100,
		MaxCargo: 10,
		Cargo: map[ResourceType]int{
			Iron: 0,
			Food: 0,
			Fuel: 0,
		},
		Offer: map[ResourceType]int{
			Iron: 5,
			Food: 5,
			Fuel: 5,
		},
		ColonizationCooldown: time.Now().Add(-time.Hour),
	}
	s.planets = []*Planet{
		{
			Resources: map[ResourceType]int{
				Iron: 10,
				Food: 10,
				Fuel: 10,
			},
			Buildings: []*Building{},
		},
		{
			Resources: map[ResourceType]int{
				Iron: 5,
				Food: 5,
				Fuel: 5,
			},
			Buildings: []*Building{},
		},
	}
}

func (s *NPCSuite) TestUpdateTradeBuySell() {
	s.npc.Credits = 100
	s.npc.MaxCargo = 10
	s.npc.Cargo[Iron] = 0
	s.planets[0].Resources[Iron] = 10
	s.npc.Offer[Iron] = 2
	s.npc.tryBuy(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.GreaterOrEqual(s.npc.Cargo[Iron], 1)
	s.LessOrEqual(s.planets[0].Resources[Iron], 10)
	s.LessOrEqual(s.npc.Credits, 100)

	s.npc.Cargo[Iron] = 5
	s.npc.Offer[Iron] = 2
	s.planets[0].Resources[Iron] = 0
	s.npc.trySell(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.LessOrEqual(s.npc.Cargo[Iron], 5)
	s.GreaterOrEqual(s.planets[0].Resources[Iron], 0)
	s.GreaterOrEqual(s.npc.Credits, 100)
}

func (s *NPCSuite) TestUpdateTradeEdgeCases() {
	s.npc.UpdateTrade([]*Planet{}, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.Equal(100, s.npc.Credits)

	s.npc.MaxCargo = 0
	s.npc.UpdateTrade(s.planets, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.Equal(100, s.npc.Credits)

	s.npc.MaxCargo = 10
	s.planets[0].Resources[Iron] = 0
	s.planets[1].Resources[Iron] = 0
	s.npc.UpdateTrade(s.planets, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.Equal(100, s.npc.Credits)

	s.npc.Credits = 0
	s.planets[0].Resources[Iron] = 10
	s.planets[1].Resources[Iron] = 5
	s.npc.UpdateTrade(s.planets, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.Equal(0, s.npc.Credits)

	s.npc.Cargo[Iron] = 0
	s.npc.UpdateTrade(s.planets, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.Equal(0, s.npc.Credits)
}

func (s *NPCSuite) TestTryBuyAllBranches() {
	s.npc.MaxCargo = 0
	s.npc.tryBuy(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.tryBuy(s.planets[1], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.MaxCargo = 10
	s.planets[0].Resources[Iron] = 0
	s.planets[1].Resources[Iron] = 0
	s.npc.tryBuy(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.tryBuy(s.planets[1], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.planets[0].Resources[Iron] = 10
	s.planets[1].Resources[Iron] = 5
	s.npc.Credits = 0
	s.npc.tryBuy(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.tryBuy(s.planets[1], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
}

func (s *NPCSuite) TestTrySellAllBranches() {
	s.npc.Cargo[Iron] = 0
	s.npc.trySell(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.trySell(s.planets[1], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.Cargo[Iron] = 5
	s.npc.trySell(s.planets[0], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.npc.trySell(s.planets[1], Iron, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
}

func (s *NPCSuite) TestUpdateTradeRandom() {
	s.npc.Cargo[Iron] = 5
	s.planets[0].Resources[Iron] = 10
	s.planets[1].Resources[Iron] = 5
	s.npc.Offer[Iron] = 2
	s.npc.MaxCargo = 10
	s.npc.Credits = 100
	s.npc.UpdateTrade(s.planets, &mockRand{seekVal: 0.4, ofVal: 1}, s.log)
	s.npc.UpdateTrade(s.planets, &mockRand{seekVal: 0.6, ofVal: 1}, s.log)
}

func (s *NPCSuite) TestRunNPCLogicColonizationCityBranch() {
	p := &Planet{Buildings: []*Building{}}
	npc := &NPC{ColonizationCooldown: time.Now().Add(-time.Hour)}
	ColonizePlanet(npc, p, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.True(IsPlanetColonized(p))
	foundCity := false
	for _, b := range p.Buildings {
		if b.Type == City {
			foundCity = true
			break
		}
	}
	s.True(foundCity)
}

func (s *NPCSuite) TestRunNPCLogicColonizationMineBranch() {
	p := &Planet{Buildings: []*Building{}}
	npc := &NPC{ColonizationCooldown: time.Now().Add(-time.Hour)}
	ColonizePlanet(npc, p, &mockRand{seekVal: 0.8, ofVal: 1}, s.log)
	s.True(IsPlanetColonized(p))
	foundMine := false
	for _, b := range p.Buildings {
		if b.Type == Mine {
			foundMine = true
			break
		}
	}
	s.True(foundMine)
}

func (s *NPCSuite) TestIsPlanetColonized() {
	p := &Planet{Buildings: []*Building{}}
	s.False(IsPlanetColonized(p))
	p.Owner = &NPC{}
	s.True(IsPlanetColonized(p))
}

func (s *NPCSuite) TestExecuteTradeExternalAndInternal() {
	npc := &NPC{
		Cargo:   map[ResourceType]int{Iron: 0},
		Offer:   map[ResourceType]int{Iron: 5},
		Credits: 100,
	}
	p := &Planet{Resources: map[ResourceType]int{Iron: 10}}
	ExecuteTrade(npc, p, s.log)
	s.Less(p.Resources[Iron], 10)
	s.GreaterOrEqual(npc.Cargo[Iron], 0)
	s.LessOrEqual(npc.Credits, 100)

	p.Owner = npc
	npc.Cargo[Iron] = 0
	p.Resources[Iron] = 10
	npc.Offer[Iron] = 10
	ExecuteTrade(npc, p, s.log)
	s.Equal(10, npc.Cargo[Iron])
	s.Equal(0, p.Resources[Iron])
}

func (s *NPCSuite) TestRunNPCLogicCooldown() {
	npc := &NPC{ColonizationCooldown: time.Now().Add(time.Hour)}
	planets := []*Planet{{Buildings: []*Building{}}}
	RunNPCLogic(npc, planets, &mockRand{seekVal: 0.5, ofVal: 1}, s.log)
	s.False(IsPlanetColonized(planets[0]))
}

func (s *NPCSuite) TestRunNPCLogicTradeBranch() {
	npc := &NPC{
		Credits:              100,
		MaxCargo:             10,
		Cargo:                map[ResourceType]int{Iron: 0, Food: 0, Fuel: 0},
		Offer:                map[ResourceType]int{Iron: 5, Food: 5, Fuel: 5},
		ColonizationCooldown: time.Now().Add(-time.Hour),
	}
	planets := []*Planet{
		{
			Resources: map[ResourceType]int{Iron: 10, Food: 10, Fuel: 10},
			Buildings: []*Building{},
		},
	}
	RunNPCLogic(npc, planets, &mockRand{seekVal: 0.2, ofVal: 0}, s.log)
}

func (s *NPCSuite) TestRunNPCLogicColonizeBranch() {
	npc := &NPC{
		Credits:              100,
		MaxCargo:             10,
		Cargo:                map[ResourceType]int{Iron: 0, Food: 0, Fuel: 0},
		Offer:                map[ResourceType]int{Iron: 5, Food: 5, Fuel: 5},
		ColonizationCooldown: time.Now().Add(-time.Hour),
	}
	planets := []*Planet{
		{
			Resources: map[ResourceType]int{Iron: 10, Food: 10, Fuel: 10},
			Buildings: []*Building{},
		},
	}
	RunNPCLogic(npc, planets, &mockRand{seekVal: 0.01, ofVal: 0}, s.log)
	s.True(IsPlanetColonized(planets[0]))
}

func (s *NPCSuite) TestMinUtility() {
	s.Equal(1, min(1, 2))
	s.Equal(2, min(2, 3))
	s.Equal(-1, min(-1, 0))
	s.Equal(0, min(0, 0))
}
