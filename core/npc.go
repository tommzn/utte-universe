package core

import (
	"time"
)

func (n *NPC) UpdateTrade(planets []*Planet, rand Random, log Log) {

	if len(planets) == 0 {
		log.Info("NPC %s: No planets available for trade.", n.Name)
		return
	}

	p := planets[rand.Of(len(planets))]
	log.Debug("NPC %s: Selected planet %s for trade.", n.Name, p.Name)

	resType := ResourceType(rand.Of(len(n.Offer)))

	if rand.Seek() < 0.5 {
		log.Debug("NPC %s: Attempting to buy %v from planet %s.", n.Name, resType, p.Name)
		n.tryBuy(p, resType, rand, log)
	} else {
		log.Debug("NPC %s: Attempting to sell %v to planet %s.", n.Name, resType, p.Name)
		n.trySell(p, resType, rand, log)
	}
}

func (n *NPC) tryBuy(p *Planet, resType ResourceType, rand Random, log Log) {

	currentCargo := 0
	for _, qty := range n.Cargo {
		currentCargo += qty
	}
	if currentCargo >= n.MaxCargo {
		log.Error("NPC %s: Cargo full, cannot buy.", n.Name)
		return
	}

	if p.Resources[resType] <= 5 {
		log.Error("NPC %s: Not enough %v on planet %s to buy.", n.Name, resType, p.Name)
		return
	}

	amount := rand.Of(5) + 1
	if p.Resources[resType] < amount {
		amount = p.Resources[resType]
	}

	price := n.Offer[resType] * amount

	if n.Credits < price {
		log.Error("NPC %s: Not enough credits to buy %d units of %v.", n.Name, amount, resType)
		return
	}

	p.Resources[resType] -= amount
	n.Cargo[resType] += amount
	n.Credits -= price
	log.Info("NPC %s bought %d units of %v from planet %s.", n.Name, amount, resType, p.Name)
}

func (n *NPC) trySell(p *Planet, resType ResourceType, rand Random, log Log) {
	if n.Cargo[resType] <= 0 {
		log.Error("NPC %s: No %v to sell.", n.Name, resType)
		return
	}

	amount := rand.Of(5) + 1
	if n.Cargo[resType] < amount {
		amount = n.Cargo[resType]
	}

	price := (n.Offer[resType] + 2) * amount

	p.Resources[resType] += amount
	n.Cargo[resType] -= amount
	n.Credits += price
	log.Info("NPC %s sold %d units of %v to planet %s.", n.Name, amount, resType, p.Name)
}

func RunNPCLogic(npc *NPC, planets []*Planet, rand Random, log Log) {
	if time.Now().Before(npc.ColonizationCooldown) {
		log.Debug("NPC %s: Colonization cooldown active.", npc.Name)
		return
	}

	candidatePlanets := []*Planet{}
	for _, p := range planets {
		if !IsPlanetColonized(p) {
			candidatePlanets = append(candidatePlanets, p)
		}
	}

	if len(candidatePlanets) > 0 && rand.Seek() < 0.05 {
		planet := candidatePlanets[rand.Of(len(candidatePlanets))]
		ColonizePlanet(npc, planet, rand, log)
		npc.ColonizationCooldown = time.Now().Add(time.Duration(rand.Of(3600)+600) * time.Second)
		log.Info("NPC %s colonized planet %s.", npc.Name, planet.Name)
		return
	}

	if len(planets) > 0 && rand.Seek() < 0.3 {
		p := planets[rand.Of(len(planets))]
		log.Debug("NPC %s: Attempting trade with planet %s.", npc.Name, p.Name)
		ExecuteTrade(npc, p, log)
	}
}

func IsPlanetColonized(p *Planet) bool {
	return p.Owner != nil
}
func ColonizePlanet(npc *NPC, p *Planet, rand Random, log Log) {
	p.Owner = npc

	if rand.Seek() < 0.7 {
		city := &Building{
			Type:       City,
			Level:      1,
			Production: map[ResourceType]int{Food: 2, Iron: 2, Fuel: 1},
			Modifiers:  map[ResourceType]float64{Food: 1.0, Iron: 1.0, Fuel: 1.0},
			BuildCost:  map[ResourceType]int{Food: 10, Iron: 10, Fuel: 5},
		}
		p.Buildings = append(p.Buildings, city)
		log.Info("NPC %s established a city on planet %s.", npc.Name, p.Name)
	} else {
		mine := &Building{
			Type:       Mine,
			Level:      1,
			Production: map[ResourceType]int{Iron: 3},
			Modifiers:  map[ResourceType]float64{Iron: 1.0},
			BuildCost:  map[ResourceType]int{Food: 5, Iron: 15},
		}
		p.Buildings = append(p.Buildings, mine)
		log.Info("NPC %s established a mine on planet %s.", npc.Name, p.Name)
	}
}
func ExecuteTrade(npc *NPC, p *Planet, log Log) {
	if p.Owner == npc {
		for res, offerAmount := range npc.Offer {
			planetAmount := p.Resources[res]
			transferAmount := min(planetAmount, offerAmount)
			if transferAmount <= 0 {
				continue
			}
			p.Resources[res] -= transferAmount
			npc.Cargo[res] += transferAmount
			log.Info("NPC %s internal transfer: %d units of %v from planet %s.", npc.Name, transferAmount, res, p.Name)
		}
		return
	}

	for res, offerAmount := range npc.Offer {
		planetAmount := p.Resources[res]
		tradeAmount := min(planetAmount, offerAmount)
		if tradeAmount <= 0 {
			continue
		}
		p.Resources[res] -= tradeAmount
		npc.Cargo[res] += tradeAmount
		price := 1
		totalCost := tradeAmount * price

		if npc.Credits >= totalCost {
			npc.Credits -= totalCost
			log.Info("NPC %s external trade: bought %d units of %v from planet %s.", npc.Name, tradeAmount, res, p.Name)
		} else {
			log.Error("NPC %s: Not enough credits for external trade.", npc.Name)
		}
	}
}

// ------------------- Utility -------------------
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
