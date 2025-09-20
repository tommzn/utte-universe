package core

import (
	"context"
	"sync"
	"time"
)

type Game struct {
	mu sync.Mutex

	config Config
	random Random
	log    Log

	Planets      []*Planet
	NPCs         []*NPC
	ActiveEvents []*Event

	planetUpdates chan []*Planet
	npcUpdates    chan []*NPC
	eventUpdates  chan []*Event
}

func NewGameService(config Config, random Random, log Log, planets []*Planet, npcs []*NPC) *Game {
	return &Game{
		config:        config,
		random:        random,
		Planets:       planets,
		NPCs:          npcs,
		log:           log,
		ActiveEvents:  []*Event{},
		planetUpdates: make(chan []*Planet, 10),
		npcUpdates:    make(chan []*NPC, 10),
		eventUpdates:  make(chan []*Event, 10),
	}
}

func (g *Game) GameLoop(ctx context.Context) {

	ticker := g.newTimer()
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			g.log.Info("Game loop canceled via context.")
			return
		case <-ticker.C:
			g.mu.Lock()
			g.log.Debug("Game tick started.")
			ProduceResources(g.Planets, g.log)
			g.ActiveEvents = MaybeTriggerEvent(g.Planets, g.ActiveEvents, g.random, g.log)
			g.ActiveEvents = UpdateEvents(g.ActiveEvents, g.log)
			for _, npc := range g.NPCs {
				RunNPCLogic(npc, g.Planets, g.random, g.log)
			}

			g.sendUpdates()
			g.log.Debug("Game tick completed.")
			g.mu.Unlock()
			g.log.Flush()
		}
	}
}

func (g *Game) sendUpdates() {
	select {
	case g.planetUpdates <- g.Planets:
	default:
	}
	select {
	case g.npcUpdates <- g.NPCs:
	default:
	}
	select {
	case g.eventUpdates <- g.ActiveEvents:
	default:
	}
}

func (g *Game) newTimer() *time.Ticker {
	return time.NewTicker(g.config.TickDuration)
}
