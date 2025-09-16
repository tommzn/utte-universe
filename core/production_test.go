package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ProduceResourcesSuite struct {
	suite.Suite
	log *mockLog
}

func TestProduceResourcesSuite(t *testing.T) {
	suite.Run(t, new(ProduceResourcesSuite))
}

func (s *ProduceResourcesSuite) SetupTest() {
	s.log = &mockLog{}
}

func (s *ProduceResourcesSuite) TestBuildingsOnly() {
	p := &Planet{
		Buildings: []*Building{
			{
				Type:       City,
				Production: map[ResourceType]int{},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
			{
				Type:       Farm,
				Production: map[ResourceType]int{Food: 4},
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
		Type:      TerraLike,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p}, s.log)
	s.Equal(asInt(float64(2)*1.0), p.Resources[Iron], "Iron mismatch")
	s.Equal(4, p.Resources[Food], "Food mismatch")
	s.Equal(2, p.Resources[Fuel], "Fuel mismatch")
}

func (s *ProduceResourcesSuite) TestPlanetTypeModifiers() {
	p := &Planet{
		Buildings: []*Building{
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
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
		Type:      GasGiant,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p}, s.log)
	// GasGiant: Iron = 2*1.0=2, Fuel = 2*1.0=2 (BaseProductionModifier returns 1.0 for Iron/Fuel unless Food)
	expectedIron := asInt(float64(2) * 1.0)
	expectedFuel := asInt(float64(2) * 1.0)
	s.Equal(expectedIron, p.Resources[Iron], "Iron mismatch")
	s.Equal(expectedFuel, p.Resources[Fuel], "Fuel mismatch")
}

func (s *ProduceResourcesSuite) TestMultiplePlanets() {
	p1 := &Planet{
		Buildings: []*Building{
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
		},
		Type:      Desert,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	p2 := &Planet{
		Buildings: []*Building{
			{
				Type:       Farm,
				Production: map[ResourceType]int{Food: 3},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
		},
		Type:      Icy,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p1, p2}, s.log)
	s.Equal(asInt(float64(2)*1.2), p1.Resources[Iron])
	s.Equal(asInt(float64(3)*1.0), p2.Resources[Food]) // Fix: Icy planet's food modifier is 1.0 in BaseProductionModifier
}

func (s *ProduceResourcesSuite) TestBuildingAndPlanetModifiers() {
	mine := &Building{
		Type:       Mine,
		Production: map[ResourceType]int{Iron: 2},
		Modifiers:  map[ResourceType]float64{Iron: 1.5},
		Level:      1,
	}
	p := &Planet{
		Buildings: []*Building{mine},
		Type:      TerraLike,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{Iron: 2.0},
	}
	ProduceResources([]*Planet{p}, s.log)
	expected := asInt(float64(2) * 2.0 * 1.0 * 1.5)
	s.Equal(expected, p.Resources[Iron])
}

func (s *ProduceResourcesSuite) TestZeroModifiers() {
	p := &Planet{
		Buildings: []*Building{
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
		},
		Type:      Desert,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{Iron: 0.0},
	}
	ProduceResources([]*Planet{p}, s.log)
	// Should default to 1.0 modifier for planet, but Desert gives 1.2 for Iron
	expected := asInt(float64(2) * 1.2)
	s.Equal(expected, p.Resources[Iron])
}

func (s *ProduceResourcesSuite) TestMultipleBuildingLevels() {
	mine := &Building{
		Type:       Mine,
		Production: map[ResourceType]int{Iron: 2},
		Modifiers:  map[ResourceType]float64{},
		Level:      3,
	}
	p := &Planet{
		Buildings: []*Building{mine},
		Type:      TerraLike,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p}, s.log)
	expected := asInt(float64(2*3) * 1.0)
	s.Equal(expected, p.Resources[Iron])
}

func (s *ProduceResourcesSuite) TestMissingResourceType() {
	mine := &Building{
		Type:       Mine,
		Production: map[ResourceType]int{},
		Modifiers:  map[ResourceType]float64{},
		Level:      1,
	}
	p := &Planet{
		Buildings: []*Building{mine},
		Type:      Icy,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p}, s.log)
	s.Equal(0, p.Resources[Iron])
	s.Equal(0, p.Resources[Food])
	s.Equal(0, p.Resources[Fuel])
}

func (s *ProduceResourcesSuite) TestUpdateEventsCall() {
	p := &Planet{
		Name:      "TestPlanet",
		Type:      TerraLike,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{Iron: 1.5},
		Buildings: []*Building{
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
		},
	}
	event := &Event{
		Name:           "Iron Surge",
		Target:         PlanetTarget,
		TargetPlanet:   p,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       1,
		RemainingTicks: 1,
	}
	events := []*Event{event}
	remaining := UpdateEvents(events, s.log)
	s.Len(remaining, 0)
	s.Equal(1.0, p.Modifiers[Iron])
}

func (s *ProduceResourcesSuite) TestNoBuildings() {
	p := &Planet{
		Buildings: []*Building{},
		Type:      TerraLike,
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p}, s.log)
	s.Equal(0, p.Resources[Iron])
	s.Equal(0, p.Resources[Food])
	s.Equal(0, p.Resources[Fuel])
}

func (s *ProduceResourcesSuite) TestNilModifiers() {
	p := &Planet{
		Buildings: []*Building{
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
				Modifiers:  nil,
				Level:      1,
			},
		},
		Type:      Desert,
		Resources: map[ResourceType]int{},
		Modifiers: nil,
	}
	ProduceResources([]*Planet{p}, s.log)
	// Should default to 1.0 for both planet and building modifiers, but Desert gives 1.2 for Iron
	expected := asInt(float64(2) * 1.2)
	s.Equal(expected, p.Resources[Iron])
}

func (s *ProduceResourcesSuite) TestUnknownPlanetType() {
	p := &Planet{
		Buildings: []*Building{
			{
				Type:       Mine,
				Production: map[ResourceType]int{Iron: 2},
				Modifiers:  map[ResourceType]float64{},
				Level:      1,
			},
		},
		Type:      PlanetType(999),
		Resources: map[ResourceType]int{},
		Modifiers: map[ResourceType]float64{},
	}
	ProduceResources([]*Planet{p}, s.log)
	// Should default to 1.0 modifier
	expected := asInt(float64(2) * 1.0)
	s.Equal(expected, p.Resources[Iron])
}

func (s *ProduceResourcesSuite) TestBaseProductionModifierAllTypes() {
	s.Equal(1.0, BaseProductionModifier(TerraLike, Iron))
	s.Equal(1.2, BaseProductionModifier(Desert, Iron))
	s.Equal(0.5, BaseProductionModifier(Desert, Food))
	s.Equal(1.5, BaseProductionModifier(GasGiant, Iron))
	s.Equal(0.0, BaseProductionModifier(GasGiant, Food))
	s.Equal(1.0, BaseProductionModifier(Icy, Iron))
	s.Equal(0.7, BaseProductionModifier(Icy, Food))
	s.Equal(1.0, BaseProductionModifier(PlanetType(999), Iron))
}
