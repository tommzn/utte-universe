package core

import (
	"context"
	"sync"
	"time"
)

type Game struct {
	config       Config
	random       Random
	Planets      []*Planet
	NPCs         []*NPC
	ActiveEvents []*Event
	log          Log
	mu           sync.Mutex
}

func NewGameService(config Config, random Random, log Log, planets []*Planet, npcs []*NPC) *Game {
	return &Game{
		config:       config,
		random:       random,
		Planets:      planets,
		NPCs:         npcs,
		log:          log,
		ActiveEvents: []*Event{},
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
			g.log.Debug("Game tick completed.")
			g.mu.Unlock()
		}
	}
}

func (g *Game) newTimer() *time.Ticker {
	return time.NewTicker(g.config.TickDuration)
}
