package core

func (p *Planet) CanBuild(b *Building, log Log) bool {
	// Planet type restrictions
	switch b.Type {
	case Farm:
		if p.Type != TerraLike && p.Type != Icy {
			log.Error("Cannot build Farm on planet type %v (%s)", p.Type, p.Name)
			return false
		}
	case Mine:
		if p.Type == GasGiant {
			log.Error("Cannot build Mine on GasGiant (%s)", p.Name)
			return false
		}
	case Refinery:
		// allowed everywhere
	case City:
		// allowed everywhere
	}

	// Check resource costs
	for res, cost := range b.BuildCost {
		if p.Resources[res] < cost {
			log.Error("Insufficient %v to build %v on planet %s", res, b.Type, p.Name)
			return false
		}
	}
	return true
}

func (p *Planet) Build(b *Building, log Log) bool {
	if !p.CanBuild(b, log) {
		return false
	}
	for res, cost := range b.BuildCost {
		p.Resources[res] -= cost
		log.Debug("Resource %v deducted by %d for building %v on planet %s", res, cost, b.Type, p.Name)
	}
	p.Buildings = append(p.Buildings, b)
	log.Info("Built %v on planet %s", b.Type, p.Name)
	return true
}
