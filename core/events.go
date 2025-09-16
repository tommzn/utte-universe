package core

const BaseEventChance = 0.05

func ShouldTriggerEvent(rand Random, chance float64, log Log) bool {
	result := rand.Seek() < chance
	log.Debug("ShouldTriggerEvent: chance=%.3f, result=%v", chance, result)
	return result
}

func MaybeTriggerEvent(planets []*Planet, activeEvents []*Event, rand Random, log Log) []*Event {

	if len(planets) == 0 {
		log.Info("No planets available for event triggering.")
		return activeEvents
	}

	// pick a random planet candidate
	p := planets[rand.Of(len(planets))]
	log.Debug("Selected planet %s for event consideration.", p.Name)

	// Adjust chance based on planet type
	chance := BaseEventChance
	switch p.Type {
	case TerraLike:
		chance *= 1.0
	case Desert:
		chance *= 1.5 // more frequent events
	case GasGiant:
		chance *= 0.8
	case Icy:
		chance *= 1.2
	}
	// 5% chance per tick to trigger a new event
	if ShouldTriggerEvent(rand, chance, log) {

		// select whether the event targets a building or the whole planet
		var target EventTarget
		var building *Building
		if len(p.Buildings) > 0 && rand.Seek() < 0.5 {
			target = BuildingTarget
			building = p.Buildings[rand.Of(len(p.Buildings))]
			log.Debug("Event targets building %v on planet %s.", building.Type, p.Name)
		} else {
			target = PlanetTarget
			log.Debug("Event targets planet %s.", p.Name)
		}

		// create event metadata
		name := ChooseEventName(p, target, rand)
		boost := ChooseEventBoost(p, target, building, rand)
		e := &Event{
			Name:           name,
			Target:         target,
			TargetPlanet:   p,
			TargetBuilding: building,
			ResourceBoost:  boost,
			Duration:       5,
			RemainingTicks: 5,
		}

		// apply boost immediately
		if e.Target == PlanetTarget && e.TargetPlanet != nil {
			if e.TargetPlanet.Modifiers == nil {
				e.TargetPlanet.Modifiers = make(map[ResourceType]float64)
			}
			for res, multiplier := range e.ResourceBoost {
				if e.TargetPlanet.Modifiers[res] == 0 {
					e.TargetPlanet.Modifiers[res] = 1.0
				}
				e.TargetPlanet.Modifiers[res] *= multiplier
				log.Info("Applied event '%s' boost %.2f to %v on planet %s.", e.Name, multiplier, res, p.Name)
			}
		} else if e.Target == BuildingTarget && e.TargetBuilding != nil {
			if e.TargetBuilding.Modifiers == nil {
				e.TargetBuilding.Modifiers = make(map[ResourceType]float64)
			}
			for res, multiplier := range e.ResourceBoost {
				if e.TargetBuilding.Modifiers[res] == 0 {
					e.TargetBuilding.Modifiers[res] = 1.0
				}
				e.TargetBuilding.Modifiers[res] *= multiplier
				log.Info("Applied event '%s' boost %.2f to %v on building %v.", e.Name, multiplier, res, e.TargetBuilding.Type)
			}
		}

		activeEvents = append(activeEvents, e)
		log.Info("Event '%s' triggered on planet %s.", e.Name, p.Name)
	}

	return activeEvents
}

func UpdateEvents(activeEvents []*Event, log Log) []*Event {
	var remaining []*Event
	for _, e := range activeEvents {
		e.RemainingTicks--
		if e.RemainingTicks > 0 {
			remaining = append(remaining, e)
			continue
		}

		if e.Target == PlanetTarget && e.TargetPlanet != nil {
			for res, multiplier := range e.ResourceBoost {
				if e.TargetPlanet.Modifiers == nil {
					continue
				}
				if e.TargetPlanet.Modifiers[res] == 0 {
					e.TargetPlanet.Modifiers[res] = 1.0
					continue
				}
				e.TargetPlanet.Modifiers[res] /= multiplier
				if e.TargetPlanet.Modifiers[res] == 0 {
					e.TargetPlanet.Modifiers[res] = 1.0
				}
				log.Info("Reverted event '%s' boost %.2f from %v on planet %s.", e.Name, multiplier, res, e.TargetPlanet.Name)
			}
		} else if e.Target == BuildingTarget && e.TargetBuilding != nil {
			for res, multiplier := range e.ResourceBoost {
				if e.TargetBuilding.Modifiers == nil {
					continue
				}
				if e.TargetBuilding.Modifiers[res] == 0 {
					e.TargetBuilding.Modifiers[res] = 1.0
					continue
				}
				e.TargetBuilding.Modifiers[res] /= multiplier
				if e.TargetBuilding.Modifiers[res] == 0 {
					e.TargetBuilding.Modifiers[res] = 1.0
				}
				log.Info("Reverted event '%s' boost %.2f from %v on building %v.", e.Name, multiplier, res, e.TargetBuilding.Type)
			}
		}
	}
	return remaining
}

func ChooseEventName(p *Planet, target EventTarget, rand Random) string {

	if target == BuildingTarget && len(p.Buildings) > 0 {
		b := p.Buildings[rand.Of(len(p.Buildings))]
		switch b.Type {
		case Mine:
			if rand.Seek() < 0.7 {
				return "Mine Collapse"
			}
			return "Iron Boom"
		case Farm:
			if p.Type == Desert && rand.Seek() < 0.5 {
				return "Drought"
			}
			return "Bountiful Harvest"
		case Refinery:
			return "Fuel Boost"
		case City:
			return "Economic Boom"
		}
	} else {
		switch p.Type {
		case Desert:
			if rand.Seek() < 0.6 {
				return "Heatwave"
			}
			return "Resource Windfall"
		case GasGiant:
			return "Storm Surge"
		case Icy:
			return "Ice Storm"
		case TerraLike:
			return "Normal Fluctuation"
		}
	}
	return "Generic Event"
}

func ChooseEventBoost(p *Planet, target EventTarget, b *Building, rand Random) map[ResourceType]float64 {

	boost := map[ResourceType]float64{
		Iron: 1.0,
		Food: 1.0,
		Fuel: 1.0,
	}

	if target == BuildingTarget && b != nil {
		// Building-targeted event: stronger effects on the building's main resource
		switch b.Type {
		case Mine:
			// Mine events: big effect on Iron
			if rand.Seek() < 0.7 {
				boost[Iron] = 1.5 // boom
			} else {
				boost[Iron] = 0.5 // partial collapse
			}
			// small side effects
			boost[Food] = 1.0
			boost[Fuel] = 1.0

		case Farm:
			if p.Type == Desert && rand.Seek() < 0.5 {
				boost[Food] = 0.5 // drought on farm
			} else {
				boost[Food] = 1.4 // bountiful harvest
			}
			boost[Iron] = 1.0
			boost[Fuel] = 1.0

		case Refinery:
			// refinery events mostly affect Fuel
			if rand.Seek() < 0.6 {
				boost[Fuel] = 1.6
			} else {
				boost[Fuel] = 0.8
			}
			boost[Iron] = 1.0
			boost[Food] = 1.0

		case City:
			// economic events - small boosts across the board
			boost[Iron] = 1.1
			boost[Food] = 1.1
			boost[Fuel] = 1.0
		}
		return boost
	}

	// Planet-targeted events: depend mostly on planet type
	switch p.Type {
	case Desert:
		// deserts: food is vulnerable
		boost[Food] = 0.5
		boost[Iron] = 1.2
		boost[Fuel] = 1.0
	case GasGiant:
		boost[Food] = 0.0
		boost[Iron] = 1.3
		boost[Fuel] = 1.5
	case Icy:
		boost[Food] = 0.7
		boost[Iron] = 1.0
		boost[Fuel] = 1.0
	case TerraLike:
		boost[Food] = 1.0
		boost[Iron] = 1.0
		boost[Fuel] = 1.0
	default:
		boost[Food] = 1.0
		boost[Iron] = 1.0
		boost[Fuel] = 1.0
	}

	// small random variation so events feel less deterministic
	if rand.Seek() < 0.1 {
		boost[Iron] *= 1.0 + (rand.Seek()-0.5)*0.1
		boost[Food] *= 1.0 + (rand.Seek()-0.5)*0.1
		boost[Fuel] *= 1.0 + (rand.Seek()-0.5)*0.1
	}

	return boost
}
