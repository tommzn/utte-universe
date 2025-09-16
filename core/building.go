package core

func (b *Building) Upgrade(p *Planet, log Log) bool {
	upgradeCost := make(map[ResourceType]int)
	for res, cost := range b.BuildCost {
		upgradeCost[res] = cost * (b.Level + 1)
	}

	for res, cost := range upgradeCost {
		if p.Resources[res] < cost {
			log.Error("Upgrade failed: insufficient resources for %v on planet %s", res, p.Name)
			return false
		}
	}

	for res, cost := range upgradeCost {
		p.Resources[res] -= cost
		log.Debug("Resource %v deducted by %d for upgrade on planet %s", res, cost, p.Name)
	}
	b.Level++
	log.Info("Building %v upgraded to level %d on planet %s", b.Type, b.Level, p.Name)
	return true
}
