package ai

import (
	"fmt"
	"math/rand"

	"github.com/rafaelrc7/inf1771-battlebot/gamemap"
)

const (
	EXPLORING     = iota
	FETCHING_GOLD = iota
	FETCHING_PU   = iota
	ATTACKING     = iota
	RUNNING       = iota
)

type AI struct {
	State        int
	ActionStack  []int
	Coord        gamemap.Coord
	Energy       int
	TimesFired   int
	TimeRunnning int
	Gamemap      *gamemap.Map
}

func AIInit(m *gamemap.Map, c gamemap.Coord) AI {
	return AI{
		Gamemap:     m,
		State:       EXPLORING,
		ActionStack: []int{},
		Coord:       c,
		Energy:      100,
	}
}

func (ai *AI) GetDecision() int {
	switch ai.State {
	case EXPLORING:
		if len(ai.ActionStack) == 0 {
			dest := FindUnexplored(ai.Gamemap, ai.Coord)
			ai.ActionStack, _ = Astar(ai.Coord, dest, ai.Gamemap)
			fmt.Printf("Going to: (%d, %d)\n", dest.X, dest.Y)
		}
		l := len(ai.ActionStack)
		action := ai.ActionStack[l-1]
		ai.ActionStack = ai.ActionStack[:l-1]
		return action
	default:
		return rand.Intn(7)
	}
}

func FindUnexplored(m *gamemap.Map, c gamemap.Coord) gamemap.Coord {
	adjs := m.GetAdjacentCells(c)

	for _, adj := range adjs {
		if !m.Cells[adj.X][adj.Y].Visited {
			return adj
		}
	}

	return c
}
