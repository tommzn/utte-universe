package core

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SeedUniverseSuite struct {
	suite.Suite
}

func TestSeedUniverseSuite(t *testing.T) {
	suite.Run(t, new(SeedUniverseSuite))
}

func (s *SeedUniverseSuite) TestSeedUniverseDeterministic() {
	seedConfig := SeedConfig{
		NumberOfPlanets: intRange{Min: 50, Max: 50},
		MPCConfig: NPCSeedConfig{
			NumberOfNPCs: intRange{Min: 20, Max: 20},
			Credits:      intRange{Min: 201, Max: 201},
			MaxCargo:     intRange{Min: 10, Max: 10},
			Offers: map[ResourceType]intRange{
				Iron: {Min: 1, Max: 1},
				Food: {Min: 1, Max: 1},
				Fuel: {Min: 1, Max: 1},
			},
			ColonizationCooldownSeconds: 3600,
		},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 51, Max: 51},
			Food: {Min: 51, Max: 51},
			Fuel: {Min: 21, Max: 21},
		},
		BuildingChance: map[BuildingType]float64{
			Mine:     1.0,
			Farm:     1.0,
			Refinery: 1.0,
		},
		Production: intRange{Min: 4, Max: 4},
		BuildCosts: map[BuildingType]map[ResourceType]int{
			Mine:     {Iron: 20},
			Farm:     {Iron: 10},
			Refinery: {Iron: 15},
		},
	}
	r := &mockRand{seekVal: 0.5, ofVal: 1}
	planets, npcs := SeedUniverse(seedConfig, r)

	s.Equal(50, len(planets), "Should create 50 planets")
	s.Equal(20, len(npcs), "Should create 20 NPCs")

	for _, planet := range planets {
		s.True(planet.Name != "", "Planet should have a name")
		s.GreaterOrEqual(len(planet.Buildings), 1)
	}

	for _, npc := range npcs {
		s.NotEmpty(npc.Name)
		s.GreaterOrEqual(npc.Credits, 200)
		s.NotNil(npc.Cargo)
	}
}

func (s *SeedUniverseSuite) TestSeedUniverseBasic() {
	seedConfig := SeedConfig{
		NumberOfPlanets: intRange{Min: 2, Max: 2},
		MPCConfig: NPCSeedConfig{
			NumberOfNPCs: intRange{Min: 2, Max: 2},
			Credits:      intRange{Min: 201, Max: 201},
			MaxCargo:     intRange{Min: 10, Max: 10},
			Offers: map[ResourceType]intRange{
				Iron: {Min: 1, Max: 1},
				Food: {Min: 1, Max: 1},
				Fuel: {Min: 1, Max: 1},
			},
			ColonizationCooldownSeconds: 3600,
		},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 51, Max: 51},
			Food: {Min: 51, Max: 51},
			Fuel: {Min: 21, Max: 21},
		},
		BuildingChance: map[BuildingType]float64{
			Mine:     1.0,
			Farm:     1.0,
			Refinery: 1.0,
		},
		Production: intRange{Min: 4, Max: 4},
		BuildCosts: map[BuildingType]map[ResourceType]int{
			Mine:     {Iron: 20},
			Farm:     {Iron: 10},
			Refinery: {Iron: 15},
		},
	}
	planets, npcs := SeedUniverse(seedConfig, &mockRand{seekVal: 0.5, ofVal: 1})

	s.Equal(2, len(planets), "Should create 2 planets")
	s.Equal(2, len(npcs), "Should create 2 NPCs")

	// Aurora-A
	s.Equal("Aurora-A", planets[0].Name)
	s.True(planets[0].Type >= TerraLike && planets[0].Type <= Icy)
	s.Equal(51, planets[0].Resources[Iron])
	s.Equal(51, planets[0].Resources[Food])
	s.Equal(21, planets[0].Resources[Fuel])
	s.GreaterOrEqual(len(planets[0].Buildings), 1)
	s.Equal(Mine, planets[0].Buildings[0].Type)

	// Vega-B
	s.Equal("Vega-B", planets[1].Name)
	s.True(planets[1].Type >= TerraLike && planets[1].Type <= Icy)
	s.Equal(51, planets[1].Resources[Iron])
	s.Equal(51, planets[1].Resources[Food])
	s.Equal(21, planets[1].Resources[Fuel])
	s.GreaterOrEqual(len(planets[1].Buildings), 1)
	s.Equal(Mine, planets[1].Buildings[0].Type)

	s.Equal("Trader Joe-A", npcs[0].Name)
	s.GreaterOrEqual(npcs[0].Credits, 201)
	s.Equal("Merchant Mia-B", npcs[1].Name)
	s.GreaterOrEqual(npcs[1].Credits, 201)
}

func (s *SeedUniverseSuite) TestPlanetBuildingProperties() {
	seedConfig := SeedConfig{
		NumberOfPlanets: intRange{Min: 2, Max: 2},
		MPCConfig: NPCSeedConfig{
			NumberOfNPCs: intRange{Min: 2, Max: 2},
			Credits:      intRange{Min: 201, Max: 201},
			MaxCargo:     intRange{Min: 10, Max: 10},
			Offers: map[ResourceType]intRange{
				Iron: {Min: 1, Max: 1},
				Food: {Min: 1, Max: 1},
				Fuel: {Min: 1, Max: 1},
			},
			ColonizationCooldownSeconds: 3600,
		},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 51, Max: 51},
			Food: {Min: 51, Max: 51},
			Fuel: {Min: 21, Max: 21},
		},
		BuildingChance: map[BuildingType]float64{
			Mine:     1.0,
			Farm:     1.0,
			Refinery: 1.0,
		},
		Production: intRange{Min: 1, Max: 2},
		BuildCosts: map[BuildingType]map[ResourceType]int{
			Mine:     {Iron: 1},
			Farm:     {Iron: 2},
			Refinery: {Iron: 2},
		},
	}
	planets, _ := SeedUniverse(seedConfig, &mockRand{seekVal: 0.5, ofVal: 1})

	// Check all planets and their buildings for expected values
	for _, planet := range planets {
		s.NotEmpty(planet.Name)
		s.True(planet.Type >= TerraLike && planet.Type <= Icy)
		s.Equal(51, planet.Resources[Iron])
		s.Equal(51, planet.Resources[Food])
		s.Equal(21, planet.Resources[Fuel])
		s.GreaterOrEqual(len(planet.Buildings), 1)

		for _, b := range planet.Buildings {
			s.Contains([]BuildingType{Mine, Farm, Refinery}, b.Type)
			s.Equal(1, b.Level)
			// Production should be between 1 and 2
			for _, v := range b.Production {
				s.True(v >= 1 && v <= 2)
			}
			for _, v := range b.Modifiers {
				s.Equal(1.0, v)
			}
			// BuildCost should match config
			for res, cost := range seedConfig.BuildCosts[b.Type] {
				s.Equal(cost, b.BuildCost[res])
			}
		}
	}
}

func (s *SeedUniverseSuite) TestSeedUniverseRandomized() {
	seedConfig := SeedConfig{
		NumberOfPlanets: intRange{Min: 5, Max: 5},
		MPCConfig: NPCSeedConfig{
			NumberOfNPCs: intRange{Min: 3, Max: 3},
			Credits:      intRange{Min: 1, Max: 1000},
			MaxCargo:     intRange{Min: 1, Max: 100},
			Offers: map[ResourceType]intRange{
				Iron: {Min: 1, Max: 100},
				Food: {Min: 1, Max: 100},
				Fuel: {Min: 1, Max: 100},
			},
			ColonizationCooldownSeconds: 3600,
		},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 1, Max: 100},
			Food: {Min: 1, Max: 100},
			Fuel: {Min: 1, Max: 100},
		},
		BuildingChance: map[BuildingType]float64{
			Mine:     1.0,
			Farm:     1.0,
			Refinery: 1.0,
		},
		Production: intRange{Min: 1, Max: 10},
		BuildCosts: map[BuildingType]map[ResourceType]int{
			Mine:     {Iron: 20},
			Farm:     {Iron: 10},
			Refinery: {Iron: 15},
		},
	}
	r := &mockRand{seekVal: 0.6, ofVal: 2}
	planets, npcs := SeedUniverse(seedConfig, r)

	s.Equal(5, len(planets), "Should create correct number of planets")
	s.Equal(3, len(npcs), "Should create correct number of NPCs")

	for _, planet := range planets {
		s.NotEmpty(planet.Name)
		s.True(planet.Type >= TerraLike && planet.Type <= Icy)
		s.GreaterOrEqual(planet.Resources[Iron], 1)
		s.GreaterOrEqual(planet.Resources[Food], 1)
		s.GreaterOrEqual(planet.Resources[Fuel], 1)
		s.NotNil(planet.Modifiers)
		s.GreaterOrEqual(len(planet.Buildings), 1)
		for _, b := range planet.Buildings {
			s.NotNil(b.Production)
			s.NotNil(b.Modifiers)
			s.NotNil(b.BuildCost)
			s.GreaterOrEqual(b.Level, 1)
		}
	}

	for _, npc := range npcs {
		s.NotEmpty(npc.Name)
		s.GreaterOrEqual(npc.Credits, 1)
		s.GreaterOrEqual(npc.MaxCargo, 1)
		s.NotNil(npc.Offer)
		s.NotNil(npc.Cargo)
		s.True(npc.ColonizationCooldown.After(npc.ColonizationCooldown.Add(-3600)))
	}
}

func (s *SeedUniverseSuite) TestPlanetNameGeneration() {
	names := map[string]bool{}
	for i := 0; i < 30; i++ {
		name := GeneratePlanetName(i)
		s.NotEmpty(name)
		names[name] = true
	}
	s.GreaterOrEqual(len(names), 8, "Should cycle through at least 8 unique names")
}

func (s *SeedUniverseSuite) TestNPCNameGeneration() {
	names := map[string]bool{}
	for i := 0; i < 15; i++ {
		name := GenerateNPCName(i)
		s.NotEmpty(name)
		names[name] = true
	}
	s.GreaterOrEqual(len(names), 5, "Should cycle through at least 5 unique names")
}

// Additional edge case: zero planets/NPCs
func (s *SeedUniverseSuite) TestZeroPlanetsAndNPCs() {
	seedConfig := SeedConfig{
		NumberOfPlanets: intRange{Min: 0, Max: 0},
		MPCConfig: NPCSeedConfig{
			NumberOfNPCs: intRange{Min: 0, Max: 0},
			Credits:      intRange{Min: 1, Max: 1},
			MaxCargo:     intRange{Min: 1, Max: 1},
			Offers: map[ResourceType]intRange{
				Iron: {Min: 1, Max: 1},
				Food: {Min: 1, Max: 1},
				Fuel: {Min: 1, Max: 1},
			},
			ColonizationCooldownSeconds: 3600,
		},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 1, Max: 1},
			Food: {Min: 1, Max: 1},
			Fuel: {Min: 1, Max: 1},
		},
		BuildingChance: map[BuildingType]float64{
			Mine:     1.0,
			Farm:     1.0,
			Refinery: 1.0,
		},
		Production: intRange{Min: 1, Max: 1},
		BuildCosts: map[BuildingType]map[ResourceType]int{
			Mine:     {Iron: 1},
			Farm:     {Iron: 1},
			Refinery: {Iron: 1},
		},
	}
	planets, npcs := SeedUniverse(seedConfig, &mockRand{seekVal: 0.5, ofVal: 1})
	s.Equal(0, len(planets))
	s.Equal(0, len(npcs))
}

// Additional edge case: negative values (should result in zero)
func (s *SeedUniverseSuite) TestNegativePlanetsAndNPCs() {
	seedConfig := SeedConfig{
		NumberOfPlanets: intRange{Min: -5, Max: -3},
		MPCConfig: NPCSeedConfig{
			NumberOfNPCs: intRange{Min: -5, Max: -3},
			Credits:      intRange{Min: 1, Max: 1},
			MaxCargo:     intRange{Min: 1, Max: 1},
			Offers: map[ResourceType]intRange{
				Iron: {Min: 1, Max: 1},
				Food: {Min: 1, Max: 1},
				Fuel: {Min: 1, Max: 1},
			},
			ColonizationCooldownSeconds: 3600,
		},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 1, Max: 1},
			Food: {Min: 1, Max: 1},
			Fuel: {Min: 1, Max: 1},
		},
		BuildingChance: map[BuildingType]float64{
			Mine:     1.0,
			Farm:     1.0,
			Refinery: 1.0,
		},
		Production: intRange{Min: 1, Max: 1},
		BuildCosts: map[BuildingType]map[ResourceType]int{
			Mine:     {Iron: 1},
			Farm:     {Iron: 1},
			Refinery: {Iron: 1},
		},
	}
	planets, npcs := SeedUniverse(seedConfig, &mockRand{seekVal: 0.5, ofVal: 1})
	s.Equal(0, len(planets))
	s.Equal(0, len(npcs))
}
