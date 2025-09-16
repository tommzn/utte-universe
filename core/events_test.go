package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type EventsSuite struct {
	suite.Suite
	planets      []*Planet
	activeEvents []*Event
	log          *mockLog
}

func TestEventsSuite(t *testing.T) {
	suite.Run(t, new(EventsSuite))
}

func (s *EventsSuite) SetupTest() {
	s.planets = []*Planet{
		{Name: "Earth", Modifiers: map[ResourceType]float64{Iron: 1.0}},
		{Name: "Mars", Modifiers: map[ResourceType]float64{Iron: 1.0}},
	}
	s.activeEvents = []*Event{}
	s.log = &mockLog{}
}

func (s *EventsSuite) TestMaybeTriggerEventNoPlanets() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	events := MaybeTriggerEvent([]*Planet{}, s.activeEvents, r, s.log)
	s.Equal(s.activeEvents, events)
}

func (s *EventsSuite) TestMaybeTriggerEventTriggers() {
	r := &mockRand{seekVal: 0.01, ofVal: 1}
	events := MaybeTriggerEvent(s.planets, s.activeEvents, r, s.log)
	s.Len(events, 1)
	event := events[0]
	s.Contains([]string{
		"Heatwave", "Resource Windfall", "Storm Surge", "Ice Storm", "Normal Fluctuation", "Generic Event",
	}, event.Name)
	s.Equal(PlanetTarget, event.Target)
	s.Contains(s.planets, event.TargetPlanet)
	s.InDelta(1.0, event.ResourceBoost[Iron], 0.5)
	s.InDelta(1.0, event.TargetPlanet.Modifiers[Iron], 0.5)
}

func (s *EventsSuite) TestMaybeTriggerEventBuildingTarget() {
	building := &Building{
		Type:      Mine,
		Modifiers: map[ResourceType]float64{Iron: 1.0},
	}
	p := &Planet{
		Name:      "Venus",
		Type:      TerraLike,
		Buildings: []*Building{building},
		Modifiers: map[ResourceType]float64{},
	}
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	events := MaybeTriggerEvent([]*Planet{p}, []*Event{}, r, s.log)
	s.Len(events, 1)
	event := events[0]
	s.Equal(BuildingTarget, event.Target)
	s.Equal(building, event.TargetBuilding)
	s.Equal(1.5, event.ResourceBoost[Iron])
	s.Equal(1.5, building.Modifiers[Iron])
}

func (s *EventsSuite) TestUpdateEventsExpiresAndReverts() {
	planet := &Planet{Name: "Earth", Modifiers: map[ResourceType]float64{Iron: 1.5}}
	event := &Event{
		Name:           "Iron Boom",
		Target:         PlanetTarget,
		TargetPlanet:   planet,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       1,
		RemainingTicks: 1,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 0)
	s.Equal(1.0, planet.Modifiers[Iron])
}

func (s *EventsSuite) TestUpdateEventsStillActive() {
	planet := &Planet{Name: "Mars", Modifiers: map[ResourceType]float64{Iron: 1.5}}
	event := &Event{
		Name:           "Iron Boom",
		Target:         PlanetTarget,
		TargetPlanet:   planet,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       2,
		RemainingTicks: 2,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 1)
	s.Equal(1.5, planet.Modifiers[Iron])
}

func (s *EventsSuite) TestUpdateEventsNilModifiers() {
	planet := &Planet{Name: "Pluto", Modifiers: nil}
	event := &Event{
		Name:           "Ice Storm",
		Target:         PlanetTarget,
		TargetPlanet:   planet,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.2},
		Duration:       1,
		RemainingTicks: 1,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 0)
	s.Nil(planet.Modifiers)
}

func (s *EventsSuite) TestShouldTriggerEventTrue() {
	r := &mockRand{seekVal: 0.01}
	s.True(ShouldTriggerEvent(r, 0.05, s.log))
}

func (s *EventsSuite) TestShouldTriggerEventFalse() {
	r := &mockRand{seekVal: 0.99}
	s.False(ShouldTriggerEvent(r, 0.05, s.log))
}

func (s *EventsSuite) TestMaybeTriggerEventNoBuildings() {
	p := &Planet{
		Name:      "NoBuildings",
		Type:      Desert,
		Buildings: []*Building{},
		Modifiers: map[ResourceType]float64{},
	}
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	events := MaybeTriggerEvent([]*Planet{p}, []*Event{}, r, s.log)
	s.Len(events, 1)
	event := events[0]
	s.Equal(PlanetTarget, event.Target)
	s.Equal(p, event.TargetPlanet)
}

func (s *EventsSuite) TestMaybeTriggerEventNilModifiers() {
	p := &Planet{
		Name:      "NilModifiers",
		Type:      Desert,
		Buildings: []*Building{},
		Modifiers: nil,
	}
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	events := MaybeTriggerEvent([]*Planet{p}, []*Event{}, r, s.log)
	s.Len(events, 1)
	event := events[0]
	s.NotNil(event.TargetPlanet.Modifiers)
}

func (s *EventsSuite) TestUpdateEventsBuildingReverts() {
	building := &Building{
		Type:      Mine,
		Modifiers: map[ResourceType]float64{Iron: 1.5},
	}
	event := &Event{
		Name:           "Mine Collapse",
		Target:         BuildingTarget,
		TargetBuilding: building,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       1,
		RemainingTicks: 1,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 0)
	s.Equal(1.0, building.Modifiers[Iron])
}

func (s *EventsSuite) TestUpdateEventsBuildingNilModifiers() {
	building := &Building{
		Type:      Mine,
		Modifiers: nil,
	}
	event := &Event{
		Name:           "Mine Collapse",
		Target:         BuildingTarget,
		TargetBuilding: building,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       1,
		RemainingTicks: 1,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 0)
	s.Nil(building.Modifiers)
}

func (s *EventsSuite) TestUpdateEventsPlanetZeroModifier() {
	planet := &Planet{Name: "ZeroMod", Modifiers: map[ResourceType]float64{Iron: 0}}
	event := &Event{
		Name:           "Iron Boom",
		Target:         PlanetTarget,
		TargetPlanet:   planet,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       1,
		RemainingTicks: 1,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 0)
	s.Equal(1.0, planet.Modifiers[Iron])
}

func (s *EventsSuite) TestUpdateEventsBuildingZeroModifier() {
	building := &Building{
		Type:      Mine,
		Modifiers: map[ResourceType]float64{Iron: 0},
	}
	event := &Event{
		Name:           "Mine Collapse",
		Target:         BuildingTarget,
		TargetBuilding: building,
		ResourceBoost:  map[ResourceType]float64{Iron: 1.5},
		Duration:       1,
		RemainingTicks: 1,
	}
	activeEvents := []*Event{event}
	remaining := UpdateEvents(activeEvents, s.log)
	s.Len(remaining, 0)
	s.Equal(1.0, building.Modifiers[Iron])
}

func (s *EventsSuite) TestChooseEventNamePlanetTypes() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: Desert, Buildings: []*Building{}}
	name := ChooseEventName(p, PlanetTarget, r)
	s.Contains([]string{"Heatwave", "Resource Windfall"}, name)

	p.Type = GasGiant
	name = ChooseEventName(p, PlanetTarget, r)
	s.Equal("Storm Surge", name)

	p.Type = Icy
	name = ChooseEventName(p, PlanetTarget, r)
	s.Equal("Ice Storm", name)

	p.Type = TerraLike
	name = ChooseEventName(p, PlanetTarget, r)
	s.Equal("Normal Fluctuation", name)
}

func (s *EventsSuite) TestChooseEventNameDefault() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: PlanetType(999)}
	name := ChooseEventName(p, PlanetTarget, r)
	s.Equal("Generic Event", name)
}

func (s *EventsSuite) TestChooseEventNameBuildingTypes() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: TerraLike}
	p.Buildings = []*Building{{Type: Mine}, {Type: Farm}, {Type: Refinery}, {Type: City}}

	name := ChooseEventName(p, BuildingTarget, r)
	s.Contains([]string{"Mine Collapse", "Iron Boom", "Bountiful Harvest", "Fuel Boost", "Economic Boom"}, name)
}

func (s *EventsSuite) TestChooseEventBoostPlanetTypes() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: Desert}
	boost := ChooseEventBoost(p, PlanetTarget, nil, r)
	s.InDelta(1.2, boost[Iron], 0.1)
	s.InDelta(0.5, boost[Food], 0.1)
	s.InDelta(1.0, boost[Fuel], 0.1)

	p.Type = GasGiant
	boost = ChooseEventBoost(p, PlanetTarget, nil, r)
	s.InDelta(1.3, boost[Iron], 0.1)
	s.InDelta(0.0, boost[Food], 0.1)
	s.InDelta(1.5, boost[Fuel], 0.1)

	p.Type = Icy
	boost = ChooseEventBoost(p, PlanetTarget, nil, r)
	s.InDelta(1.0, boost[Iron], 0.1)
	s.InDelta(0.7, boost[Food], 0.1)
	s.InDelta(1.0, boost[Fuel], 0.1)

	p.Type = TerraLike
	boost = ChooseEventBoost(p, PlanetTarget, nil, r)
	s.InDelta(1.0, boost[Iron], 0.1)
	s.InDelta(1.0, boost[Food], 0.1)
	s.InDelta(1.0, boost[Fuel], 0.1)
}

func (s *EventsSuite) TestChooseEventBoostDefaultType() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: PlanetType(999)}
	boost := ChooseEventBoost(p, PlanetTarget, nil, r)
	s.InDelta(1.0, boost[Iron], 0.1)
	s.InDelta(1.0, boost[Food], 0.1)
	s.InDelta(1.0, boost[Fuel], 0.1)
}

func (s *EventsSuite) TestChooseEventBoostBuildingTypes() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: TerraLike}
	mine := &Building{Type: Mine}
	farm := &Building{Type: Farm}
	refinery := &Building{Type: Refinery}
	city := &Building{Type: City}

	boost := ChooseEventBoost(p, BuildingTarget, mine, r)
	s.Equal(1.5, boost[Iron])
	s.Equal(1.0, boost[Food])
	s.Equal(1.0, boost[Fuel])

	boost = ChooseEventBoost(p, BuildingTarget, farm, r)
	s.Equal(1.0, boost[Iron])
	s.Equal(1.4, boost[Food])
	s.Equal(1.0, boost[Fuel])

	boost = ChooseEventBoost(p, BuildingTarget, refinery, r)
	s.Equal(1.0, boost[Iron])
	s.Equal(1.0, boost[Food])
	s.Equal(1.6, boost[Fuel])

	boost = ChooseEventBoost(p, BuildingTarget, city, r)
	s.Equal(1.1, boost[Iron])
	s.Equal(1.1, boost[Food])
	s.Equal(1.0, boost[Fuel])
}

func (s *EventsSuite) TestChooseEventBoostRandomVariation() {
	r := &mockRand{seekVal: 0.09, ofVal: 0}
	p := &Planet{Type: TerraLike}
	boost := ChooseEventBoost(p, PlanetTarget, nil, r)
	s.InDelta(1.0, boost[Iron], 0.1)
	s.InDelta(1.0, boost[Food], 0.1)
	s.InDelta(1.0, boost[Fuel], 0.1)
}

func (s *EventsSuite) TestChooseEventBoostBuildingUnknownType() {
	r := &mockRand{seekVal: 0.01, ofVal: 0}
	p := &Planet{Type: TerraLike}
	b := &Building{Type: BuildingType(999)}
	boost := ChooseEventBoost(p, BuildingTarget, b, r)
	s.Equal(1.0, boost[Iron])
	s.Equal(1.0, boost[Food])
	s.Equal(1.0, boost[Fuel])
}

func (s *EventsSuite) TestShouldTriggerEventEdgeCases() {
	r := &mockRand{seekVal: 0.0}
	s.False(ShouldTriggerEvent(r, 0.0, s.log))
	r.seekVal = 1.0
	s.False(ShouldTriggerEvent(r, 0.0, s.log))
}
