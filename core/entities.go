package core

import "time"

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

type Building struct {
	Type       BuildingType
	Level      int
	Production map[ResourceType]int     // output per tick
	Modifiers  map[ResourceType]float64 // multipliers
	BuildCost  map[ResourceType]int     // cost to build/upgrade
}

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

type Planet struct {
	Name      string                   `json:"name"`
	Type      PlanetType               `json:"type"`
	Resources map[ResourceType]int     `json:"resources"`
	Modifiers map[ResourceType]float64 `json:"modifiers"`
	Buildings []*Building              `json:"buildings"`
	Owner     *NPC                     `json:"owner"`
}

type NPC struct {
	Name                 string               `json:"name"`
	Offer                map[ResourceType]int `json:"offer"`
	Credits              int                  `json:"credits"`
	Cargo                map[ResourceType]int `json:"cargo"`
	MaxCargo             int                  `json:"maxCargo"`
	ColonizationCooldown time.Time            `json:"colonizationCooldown"`
}

type TradeAction struct {
	NPCName    string
	PlanetName string
	Resource   ResourceType
	Amount     int
}

type EventTarget int

const (
	PlanetTarget EventTarget = iota
	BuildingTarget
)

type Event struct {
	Name           string
	Target         EventTarget
	TargetPlanet   *Planet
	TargetBuilding *Building
	ResourceBoost  map[ResourceType]float64 // multiplier applied
	Duration       int                      // total ticks
	RemainingTicks int                      // ticks left
}
