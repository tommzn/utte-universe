package core

import (
	"time"
)

var (
	planetTypes              = []PlanetType{TerraLike, Desert, GasGiant, Icy}
	resourceTypes            = []ResourceType{Iron, Food, Fuel}
	buildingTypes            = []BuildingType{Mine, Farm, Refinery}
	buildingResourceRelation = map[BuildingType]ResourceType{Mine: Iron, Farm: Food, Refinery: Fuel}
	baseModifiers            = map[ResourceType]float64{Iron: 1.0, Food: 1.0, Fuel: 1.0}
)

func SeedUniverse(seedConfig SeedConfig, rand Random) ([]*Planet, []*NPC) {
	return GeneratePlanets(seedConfig, rand), GenerateNPCs(seedConfig, rand)
}

func GeneratePlanetName(idx int) string {
	names := []string{"Aurora", "Vega", "Nova", "Luna", "Terra", "Ceres", "Orion", "Eos"}
	return names[idx%len(names)] + "-" + string('A'+rune(idx%26))
}

func GenerateNPCName(idx int) string {
	names := []string{"Trader Joe", "Merchant Mia", "Captain Rex", "Baroness Lila", "Drake"}
	return names[idx%len(names)] + "-" + string('A'+rune(idx%26))
}

func GeneratePlanets(seedConfig SeedConfig, rand Random) []*Planet {

	numPlanets := rand.OfIntRange(seedConfig.NumberOfPlanets)
	planets := make([]*Planet, 0)
	for i := 0; i < numPlanets; i++ {

		planetType := PlanetType(rand.Of(len(planetTypes)))
		planets = append(planets, &Planet{
			Name:      GeneratePlanetName(i),
			Type:      planetType,
			Resources: GenerateResources(seedConfig, rand),
			Modifiers: baseModifiers,
			Buildings: GenerateBuildings(planetType, seedConfig, rand),
		})
	}
	return planets
}

func GenerateBuildings(planetType PlanetType, seedConfig SeedConfig, rand Random) []*Building {

	buildings := make([]*Building, 0)
	for _, buildingType := range buildingTypes {

		if buildingType == Farm && planetType != TerraLike && planetType != Icy {
			continue
		}

		if rand.Seek() < seedConfig.BuildingChance[buildingType] {
			resourceType := buildingResourceRelation[buildingType]
			buildings = append(buildings, &Building{
				Type:  buildingType,
				Level: 1,
				Production: map[ResourceType]int{
					resourceType: rand.OfIntRange(seedConfig.Production),
				},
				Modifiers: map[ResourceType]float64{resourceType: 1.0},
				BuildCost: seedConfig.BuildCosts[buildingType],
			})
		}
	}
	return buildings
}

func GenerateResources(seedConfig SeedConfig, rand Random) map[ResourceType]int {

	resources := make(map[ResourceType]int)
	for _, resource := range resourceTypes {
		resources[resource] = rand.OfIntRange(seedConfig.Resources[resource])
	}
	return resources
}

func GenerateNPCs(seedConfig SeedConfig, rand Random) []*NPC {

	numNPCs := rand.OfIntRange(seedConfig.MPCConfig.NumberOfNPCs)
	npcs := make([]*NPC, 0)
	for i := 0; i < numNPCs; i++ {

		offer := make(map[ResourceType]int)
		for _, resourceType := range resourceTypes {
			offer[resourceType] = rand.OfRange(
				seedConfig.MPCConfig.Offers[resourceType].Min,
				seedConfig.MPCConfig.Offers[resourceType].Max)
		}

		cargo := map[ResourceType]int{
			Iron: 0,
			Food: 0,
			Fuel: 0,
		}
		npcs = append(npcs, &NPC{
			Name:                 GenerateNPCName(i),
			Offer:                offer,
			Credits:              rand.OfIntRange(seedConfig.MPCConfig.Credits),
			Cargo:                cargo,
			MaxCargo:             rand.OfIntRange(seedConfig.MPCConfig.MaxCargo),
			ColonizationCooldown: time.Now().Add(time.Duration(rand.Of(seedConfig.MPCConfig.ColonizationCooldownSeconds)) * time.Second),
		})
	}
	return npcs
}
