package core

import (
	"time"

	"github.com/tommzn/go-config"
)

type Config struct {
	TickDuration time.Duration
	SeedConfig   SeedConfig
}

type SeedConfig struct {
	NumberOfPlanets intRange
	Resources       map[ResourceType]intRange
	BuildingChance  map[BuildingType]float64
	BuildCosts      map[BuildingType]map[ResourceType]int
	Production      intRange
	MPCConfig       NPCSeedConfig
}

type NPCSeedConfig struct {
	NumberOfNPCs                intRange
	Offers                      map[ResourceType]intRange
	Credits                     intRange
	MaxCargo                    intRange
	ColonizationCooldownSeconds int
}

func DefaultConfig() Config {
	return Config{
		TickDuration: 2 * time.Second,
		SeedConfig:   DefaultSeedConfig(),
	}
}

type intRange struct {
	Min int
	Max int
}

type RawConfig struct {
	TickDuration string        `mapstructure:"tick_duration"`
	SeedConfig   RawSeedConfig `mapstructure:"universe_seed"`
}

type RawSeedConfig struct {
	NumberOfPlanets RawIntRange            `mapstructure:"number_of_planets"`
	Resources       []ResourceConfig       `mapstructure:"resources"`
	BuildingChance  []BuildingChanceConfig `mapstructure:"building_chance"`
	BuildCosts      []BuildCostConfig      `mapstructure:"build_costs"`
	Production      RawIntRange            `mapstructure:"production"`
	NPC             RawNPCSeedConfig       `mapstructure:"npc"`
}

type ResourceConfig struct {
	Resource string `mapstructure:"resource"`
	Min      int    `mapstructure:"min"`
	Max      int    `mapstructure:"max"`
}

type BuildingChanceConfig struct {
	BuildingType string  `mapstructure:"building_type"`
	Chance       float64 `mapstructure:"chance"`
}

type BuildCostConfig struct {
	BuildingType string           `mapstructure:"building_type"`
	Resources    []ResourceAmount `mapstructure:"resources"`
}
type ResourceAmount struct {
	Resource string `mapstructure:"resource"`
	Amount   int    `mapstructure:"amount"`
}
type RawIntRange struct {
	Min int `mapstructure:"min"`
	Max int `mapstructure:"max"`
}
type RawNPCSeedConfig struct {
	NumberOfNPCs                RawIntRange      `mapstructure:"number_of_npcs"`
	Offers                      []ResourceConfig `mapstructure:"offers"`
	Credits                     RawIntRange      `mapstructure:"credits"`
	MaxCargo                    RawIntRange      `mapstructure:"max_cargo"`
	ColonizationCooldownSeconds int              `mapstructure:"colonization_cooldown_seconds"`
}

func DefaultSeedConfig() SeedConfig {
	return SeedConfig{
		NumberOfPlanets: intRange{Min: 3, Max: 15},
		Resources: map[ResourceType]intRange{
			Iron: {Min: 100, Max: 5000},
			Food: {Min: 300, Max: 3000},
			Fuel: {Min: 200, Max: 2000},
		},
		BuildingChance: DefaultBuildingChance(),
		BuildCosts:     DefaultBuildCost(),
		Production:     intRange{Min: 3, Max: 20},
		MPCConfig:      DefaultNPCSeedConfig(),
	}
}

func DefaultBuildingChance() map[BuildingType]float64 {
	return map[BuildingType]float64{
		City:     0.2,
		Mine:     0.8,
		Farm:     0.7,
		Refinery: 0.5,
	}
}

func DefaultBuildCost() map[BuildingType]map[ResourceType]int {
	return map[BuildingType]map[ResourceType]int{
		City: map[ResourceType]int{
			Iron: 100,
			Food: 50,
			Fuel: 20,
		},
		Mine: map[ResourceType]int{
			Iron: 50,
			Food: 20,
		},
		Farm: map[ResourceType]int{
			Iron: 30,
			Food: 10,
		},
		Refinery: map[ResourceType]int{
			Iron: 70,
			Fuel: 30,
		},
	}
}

func DefaultNPCSeedConfig() NPCSeedConfig {
	return NPCSeedConfig{
		NumberOfNPCs: intRange{Min: 3, Max: 8},
		Offers: map[ResourceType]intRange{
			Iron: {Min: 5, Max: 20},
			Food: {Min: 5, Max: 15},
			Fuel: {Min: 10, Max: 25},
		},
		Credits:                     intRange{Min: 200, Max: 50000},
		MaxCargo:                    intRange{Min: 50, Max: 600},
		ColonizationCooldownSeconds: 3600,
	}
}

func (c *Config) LoadFrom(conf config.Config) error {

	rawConfig := RawConfig{}
	if err := conf.Unmarshal(&rawConfig); err != nil {
		return err
	}

	// TickDuration
	if rawConfig.TickDuration != "" {
		if d, err := time.ParseDuration(rawConfig.TickDuration); err == nil {
			c.TickDuration = d
		}
	}

	// SeedConfig
	seed := rawConfig.SeedConfig

	// NumberOfPlanets
	c.SeedConfig.NumberOfPlanets = intRange{
		Min: seed.NumberOfPlanets.Min,
		Max: seed.NumberOfPlanets.Max,
	}

	// Resources
	c.SeedConfig.Resources = make(map[ResourceType]intRange)
	for _, rc := range seed.Resources {
		rt := ResourceTypeFromString(rc.Resource)
		c.SeedConfig.Resources[rt] = intRange{Min: rc.Min, Max: rc.Max}
	}

	// BuildingChance
	c.SeedConfig.BuildingChance = make(map[BuildingType]float64)
	for _, bc := range seed.BuildingChance {
		bt := BuildingTypeFromString(bc.BuildingType)
		c.SeedConfig.BuildingChance[bt] = bc.Chance
	}

	// BuildCosts
	c.SeedConfig.BuildCosts = make(map[BuildingType]map[ResourceType]int)
	for _, bc := range seed.BuildCosts {
		bt := BuildingTypeFromString(bc.BuildingType)
		resMap := make(map[ResourceType]int)
		for _, ra := range bc.Resources {
			rt := ResourceTypeFromString(ra.Resource)
			resMap[rt] = ra.Amount
		}
		c.SeedConfig.BuildCosts[bt] = resMap
	}

	// Production
	c.SeedConfig.Production = intRange{
		Min: seed.Production.Min,
		Max: seed.Production.Max,
	}

	// NPCConfig
	npc := seed.NPC
	c.SeedConfig.MPCConfig.NumberOfNPCs = intRange{Min: npc.NumberOfNPCs.Min, Max: npc.NumberOfNPCs.Max}
	c.SeedConfig.MPCConfig.Offers = make(map[ResourceType]intRange)
	for _, offer := range npc.Offers {
		rt := ResourceTypeFromString(offer.Resource)
		c.SeedConfig.MPCConfig.Offers[rt] = intRange{Min: offer.Min, Max: offer.Max}
	}
	c.SeedConfig.MPCConfig.Credits = intRange{Min: npc.Credits.Min, Max: npc.Credits.Max}
	c.SeedConfig.MPCConfig.MaxCargo = intRange{Min: npc.MaxCargo.Min, Max: npc.MaxCargo.Max}
	c.SeedConfig.MPCConfig.ColonizationCooldownSeconds = npc.ColonizationCooldownSeconds

	return nil
}
