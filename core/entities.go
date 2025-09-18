// Package core provides the main entities and types for the UTTE Universe simulation.
// It defines planets, buildings, resources, NPCs, trade actions, and events.

package core

import "time"

// PlanetType represents the type of a planet in the universe.
type PlanetType int

const (
	TerraLike PlanetType = iota
	Desert
	GasGiant
	Icy
)

func (pt PlanetType) String() string {
	switch pt {
	case TerraLike:
		return "Terra-like"
	case Desert:
		return "Desert"
	case GasGiant:
		return "Gas Giant"
	case Icy:
		return "Icy"
	default:
		return "Unknown"
	}
}

// BuildingType represents the type of a building.
type BuildingType int

const (
	Mine BuildingType = iota
	Farm
	Refinery
	City
)

func (b BuildingType) String() string {
	switch b {
	case Mine:
		return "Mine"
	case Farm:
		return "Farm"
	case Refinery:
		return "Refinery"
	case City:
		return "City"
	default:
		return "Unknown"
	}
}

// BuildingTypeFromString converts a string to a BuildingType.
func BuildingTypeFromString(s string) BuildingType {
	switch s {
	case "Mine":
		return Mine
	case "Farm":
		return Farm
	case "Refinery":
		return Refinery
	case "City":
		return City
	default:
		return BuildingType(-1) // Unknown
	}
}

// Building represents a building on a planet, including its type, level, production, modifiers, and build cost.
type Building struct {
	Type       BuildingType
	Level      int
	Production map[ResourceType]int     // output per tick
	Modifiers  map[ResourceType]float64 // multipliers
	BuildCost  map[ResourceType]int     // cost to build/upgrade
}

// ResourceType represents a type of resource in the universe.
type ResourceType int

const (
	Iron ResourceType = iota
	Food
	Fuel
)

// Helper for readable names (useful for UI/debug)
func (r ResourceType) String() string {
	switch r {
	case Iron:
		return "Iron"
	case Food:
		return "Food"
	case Fuel:
		return "Fuel"
	default:
		return "Unknown"
	}
}

// ResourceTypeFromString converts a string to a ResourceType.
func ResourceTypeFromString(s string) ResourceType {
	switch s {
	case "Iron":
		return Iron
	case "Food":
		return Food
	case "Fuel":
		return Fuel
	default:
		return ResourceType(-1) // Unknown
	}
}

// Planet represents a planet in the universe, including its type, resources, modifiers, buildings, and owner.
type Planet struct {
	Name      string                   `json:"name"`
	Type      PlanetType               `json:"type"`
	Resources map[ResourceType]int     `json:"resources"`
	Modifiers map[ResourceType]float64 `json:"modifiers"`
	Buildings []*Building              `json:"buildings"`
	Owner     *NPC                     `json:"owner"`
}

// NPC represents a non-player character, including trading offers, credits, cargo, and cooldowns.
type NPC struct {
	Name                 string               `json:"name"`
	Offer                map[ResourceType]int `json:"offer"`
	Credits              int                  `json:"credits"`
	Cargo                map[ResourceType]int `json:"cargo"`
	MaxCargo             int                  `json:"maxCargo"`
	ColonizationCooldown time.Time            `json:"colonizationCooldown"`
}

// TradeAction represents a trade action between an NPC and a planet.
type TradeAction struct {
	NPCName    string
	PlanetName string
	Resource   ResourceType
	Amount     int
}

// EventTarget specifies the target type for an event (planet or building).
type EventTarget int

const (
	PlanetTarget EventTarget = iota
	BuildingTarget
)

// Event represents a game event, which can target a planet or building and apply resource boosts for a duration.
type Event struct {
	Name           string
	Target         EventTarget
	TargetPlanet   *Planet
	TargetBuilding *Building
	ResourceBoost  map[ResourceType]float64 // multiplier applied
	Duration       int                      // total ticks
	RemainingTicks int                      // ticks left
}
