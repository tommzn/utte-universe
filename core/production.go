package core

func ProduceResources(planets []*Planet, log Log) {
	for _, p := range planets {
		for _, b := range p.Buildings {
			for resType, base := range b.Production {
				planetBoost := p.Modifiers[resType] * BaseProductionModifier(p.Type, resType)
				buildingBoost := b.Modifiers[resType]
				if planetBoost == 0 {
					planetBoost = 1.0
				}
				if buildingBoost == 0 {
					buildingBoost = 1.0
				}
				totalBoost := planetBoost * buildingBoost
				produced := int(float64(base*b.Level) * totalBoost)
				p.Resources[resType] += produced
				if produced > 0 {
					log.Info("Produced %d units of %v on planet %s (building %v, level %d)", produced, resType, p.Name, b.Type, b.Level)
				} else {
					log.Debug("No production for %v on planet %s (building %v, level %d)", resType, p.Name, b.Type, b.Level)
				}
			}
		}
	}
}

func BaseProductionModifier(pt PlanetType, res ResourceType) float64 {
	switch pt {
	case TerraLike:
		return 1.0
	case Desert:
		if res == Food {
			return 0.5
		}
		return 1.2
	case GasGiant:
		if res == Food {
			return 0.0
		}
		return 1.5
	case Icy:
		if res == Food {
			return 0.7
		}
		return 1.0
	}
	return 1.0
}
