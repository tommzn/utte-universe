package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestDefaultConfig() {
	cfg := DefaultConfig()
	s.NotNil(cfg)
	s.Equal(2*time.Second, cfg.TickDuration)
	s.Equal(DefaultSeedConfig(), cfg.SeedConfig)
}

func (s *ConfigSuite) TestDefaultSeedConfig() {
	seed := DefaultSeedConfig()
	s.Equal(intRange{Min: 3, Max: 15}, seed.NumberOfPlanets)
	s.Equal(intRange{Min: 3, Max: 8}, seed.MPCConfig.NumberOfNPCs)
	s.Equal(intRange{Min: 200, Max: 50000}, seed.MPCConfig.Credits)
	s.Equal(intRange{Min: 50, Max: 600}, seed.MPCConfig.MaxCargo)
	s.Equal(3600, seed.MPCConfig.ColonizationCooldownSeconds)
	s.Equal(intRange{Min: 3, Max: 20}, seed.Production)
	s.NotNil(seed.Resources)
	s.NotNil(seed.BuildingChance)
	s.NotNil(seed.BuildCosts)
}

func (s *ConfigSuite) TestDefaultBuildingChance() {
	chance := DefaultBuildingChance()
	s.Equal(0.2, chance[City])
	s.Equal(0.8, chance[Mine])
	s.Equal(0.7, chance[Farm])
	s.Equal(0.5, chance[Refinery])
}

func (s *ConfigSuite) TestDefaultBuildCost() {
	cost := DefaultBuildCost()
	s.Equal(100, cost[City][Iron])
	s.Equal(50, cost[City][Food])
	s.Equal(20, cost[City][Fuel])
	s.Equal(50, cost[Mine][Iron])
	s.Equal(20, cost[Mine][Food])
	s.Equal(30, cost[Farm][Iron])
	s.Equal(10, cost[Farm][Food])
	s.Equal(70, cost[Refinery][Iron])
	s.Equal(30, cost[Refinery][Fuel])
}

func (s *ConfigSuite) TestDefaultNPCSeedConfig() {
	npc := DefaultNPCSeedConfig()
	s.Equal(intRange{Min: 3, Max: 8}, npc.NumberOfNPCs)
	s.Equal(intRange{Min: 5, Max: 20}, npc.Offers[Iron])
	s.Equal(intRange{Min: 5, Max: 15}, npc.Offers[Food])
	s.Equal(intRange{Min: 10, Max: 25}, npc.Offers[Fuel])
	s.Equal(intRange{Min: 200, Max: 50000}, npc.Credits)
	s.Equal(intRange{Min: 50, Max: 600}, npc.MaxCargo)
	s.Equal(3600, npc.ColonizationCooldownSeconds)
}

func (s *ConfigSuite) TestLoadFrom() {

	cfg := &Config{}
	conf := loadConfigForTest(nil)
	err := cfg.LoadFrom(conf)
	s.NoError(err)

	// TickDuration should be overridden and positive
	s.NotEqual(DefaultConfig().TickDuration, cfg.TickDuration, "TickDuration should be overridden")
	s.Greater(cfg.TickDuration, time.Duration(0), "TickDuration should be positive")

	// NumberOfPlanets should be overridden and valid
	s.NotEqual(DefaultConfig().SeedConfig.NumberOfPlanets, cfg.SeedConfig.NumberOfPlanets, "NumberOfPlanets should be overridden")
	s.Greater(cfg.SeedConfig.NumberOfPlanets.Min, 0, "NumberOfPlanets.Min should be > 0")
	s.GreaterOrEqual(cfg.SeedConfig.NumberOfPlanets.Max, cfg.SeedConfig.NumberOfPlanets.Min, "NumberOfPlanets.Max should be >= Min")

	// Resources should not be empty and contain expected keys
	s.NotEmpty(cfg.SeedConfig.Resources, "Resources should not be empty")
	for k, v := range cfg.SeedConfig.Resources {
		s.GreaterOrEqual(v.Max, v.Min, "Resource %v: Max should be >= Min", k)
	}

	// BuildingChance should not be empty and values should be in [0,1]
	s.NotEmpty(cfg.SeedConfig.BuildingChance, "BuildingChance should not be empty")
	for k, v := range cfg.SeedConfig.BuildingChance {
		s.GreaterOrEqual(v, 0.0, "BuildingChance %v: should be >= 0", k)
		s.LessOrEqual(v, 1.0, "BuildingChance %v: should be <= 1", k)
	}

	// BuildCosts should not be empty and contain expected structure
	s.NotEmpty(cfg.SeedConfig.BuildCosts, "BuildCosts should not be empty")
	for bType, costs := range cfg.SeedConfig.BuildCosts {
		s.NotEmpty(costs, "BuildCosts for %v should not be empty", bType)
		for rType, amount := range costs {
			s.GreaterOrEqual(amount, 0, "BuildCost for %v/%v should be >= 0", bType, rType)
		}
	}

	// NPC Offers should not be empty and valid
	s.NotEmpty(cfg.SeedConfig.MPCConfig.Offers, "NPC Offers should not be empty")
	for k, v := range cfg.SeedConfig.MPCConfig.Offers {
		s.GreaterOrEqual(v.Max, v.Min, "NPC Offer %v: Max should be >= Min", k)
	}
}
