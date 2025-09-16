package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type EntitiesSuite struct {
	suite.Suite
}

func TestEntitiesSuite(t *testing.T) {
	suite.Run(t, new(EntitiesSuite))
}

func (s *EntitiesSuite) TestPlanetTypeString() {
	s.Equal("Terra-like", TerraLike.String())
	s.Equal("Desert", Desert.String())
	s.Equal("Gas Giant", GasGiant.String())
	s.Equal("Icy", Icy.String())
	s.Equal("Unknown", PlanetType(999).String())
}

func (s *EntitiesSuite) TestBuildingTypeString() {
	s.Equal("Mine", Mine.String())
	s.Equal("Farm", Farm.String())
	s.Equal("Refinery", Refinery.String())
	s.Equal("City", City.String())
	s.Equal("Unknown", BuildingType(999).String())
}

func (s *EntitiesSuite) TestResourceTypeString() {
	s.Equal("Iron", Iron.String())
	s.Equal("Food", Food.String())
	s.Equal("Fuel", Fuel.String())
	s.Equal("Unknown", ResourceType(999).String())
}

func (s *EntitiesSuite) TestBuildingStruct() {
	b := &Building{
		Type:       Mine,
		Level:      2,
		Production: map[ResourceType]int{Iron: 5},
		Modifiers:  map[ResourceType]float64{Iron: 1.2},
		BuildCost:  map[ResourceType]int{Iron: 10},
	}
	s.Equal(Mine, b.Type)
	s.Equal(2, b.Level)
	s.Equal(5, b.Production[Iron])
	s.Equal(1.2, b.Modifiers[Iron])
	s.Equal(10, b.BuildCost[Iron])
}

func (s *EntitiesSuite) TestPlanetStruct() {
	p := &Planet{
		Name:      "Earth",
		Type:      TerraLike,
		Resources: map[ResourceType]int{Iron: 100, Food: 50, Fuel: 20},
		Modifiers: map[ResourceType]float64{Iron: 1.0},
		Buildings: []*Building{
			{Type: Mine, Level: 1},
		},
	}
	s.Equal("Earth", p.Name)
	s.Equal(TerraLike, p.Type)
	s.Equal(100, p.Resources[Iron])
	s.Equal(1.0, p.Modifiers[Iron])
	s.Equal(Mine, p.Buildings[0].Type)
}

func (s *EntitiesSuite) TestNPCStruct() {
	npc := &NPC{
		Name:                 "Trader",
		Offer:                map[ResourceType]int{Iron: 5},
		Credits:              1000,
		Cargo:                map[ResourceType]int{Iron: 10},
		MaxCargo:             50,
		ColonizationCooldown: time.Now().Add(time.Hour),
	}
	s.Equal("Trader", npc.Name)
	s.Equal(5, npc.Offer[Iron])
	s.Equal(1000, npc.Credits)
	s.Equal(10, npc.Cargo[Iron])
	s.Equal(50, npc.MaxCargo)
	s.True(npc.ColonizationCooldown.After(time.Now()))
}

func (s *EntitiesSuite) TestTradeActionStruct() {
	ta := TradeAction{
		NPCName:    "Trader",
		PlanetName: "Earth",
		Resource:   Iron,
		Amount:     10,
	}
	s.Equal("Trader", ta.NPCName)
	s.Equal("Earth", ta.PlanetName)
	s.Equal(Iron, ta.Resource)
	s.Equal(10, ta.Amount)
}

func (s *EntitiesSuite) TestEventStruct() {
	p := &Planet{Name: "Earth"}
	b := &Building{Type: Mine}
	e := &Event{
		Name:           "Iron Boom",
		Target:         PlanetTarget,
		TargetPlanet:   p,
		TargetBuilding: b,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       5,
		RemainingTicks: 3,
	}
	s.Equal("Iron Boom", e.Name)
	s.Equal(PlanetTarget, e.Target)
	s.Equal(p, e.TargetPlanet)
	s.Equal(b, e.TargetBuilding)
	s.Equal(1.5, e.ResourceBoost[Iron])
	s.Equal(5, e.Duration)
	s.Equal(3, e.RemainingTicks)
}
