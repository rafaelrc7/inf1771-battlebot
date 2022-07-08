package ai

import (
	"fmt"
	"math/rand"

	"github.com/rafaelrc7/inf1771-battlebot/gamemap"
)

const (
	STOP = iota
	EXPLORING
	FETCHING_GOLD
	FETCHING_PU
	ATTACKING
	RUNNING
)

type AI struct {
	State         int
	ActionStack   []int
	Coord         gamemap.Coord
	Energy        int
	TimeRunnning  int
	Observations  uint
	EnemyDetected bool
	Gamemap       *gamemap.Map
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

func (ai *AI) GetDecision(mapChanged bool) int {
	if mapChanged {
		ai.ActionStack = []int{}
	}

	if ai.State == ATTACKING && !ai.EnemyDetected {
		ai.State = EXPLORING
	} else if ai.State != ATTACKING && ai.EnemyDetected {
		ai.State = ATTACKING
	}

	switch ai.State {
	case STOP:
		return NOTHING
	case ATTACKING:
		return ATTACK
	case EXPLORING:
		if ai.Observations&(gamemap.BLUELIGHT) != 0 {
			*ai.Gamemap.GoldCells[gamemap.CellC{X: ai.Coord.X, Y: ai.Coord.Y}] = 1500
			return TAKE
		}
		if len(ai.ActionStack) == 0 {
			dest := FindUnexplored(ai.Gamemap, ai.Coord)
			ai.ActionStack, _ = Astar(ai.Coord, dest, ai.Gamemap)
			fmt.Printf("Going to: (%d, %d)\n", dest.X, dest.Y)
		}
		if len(ai.ActionStack) == 0 {
			return NOTHING
		}

		l := len(ai.ActionStack)
		action := ai.ActionStack[l-1]
		ai.ActionStack = ai.ActionStack[:l-1]
		return action
	default:
		return NOTHING
	}
}

func FindUnexplored(m *gamemap.Map, c gamemap.Coord) gamemap.Coord {
	adjs := m.GetAdjacentCells(c)

	for _, adj := range adjs {
		if !m.Cells[adj.X][adj.Y].Visited &&
			m.Cells[adj.X][adj.Y].Status != gamemap.WALL &&
			m.Cells[adj.X][adj.Y].Status != gamemap.DANGEROUS &&
			m.Cells[adj.X][adj.Y].Status != gamemap.HOLE &&
			m.Cells[adj.X][adj.Y].Status != gamemap.TELEPORT {
			return adj
		}
	}

	for key, val := range m.GoldCells {
		if val != nil {
			if *val == 0 {
				return gamemap.Coord{X: key.X, Y: key.Y}
			}
		}
	}

	return gamemap.Coord{X: rand.Intn(m.Width + 1), Y: rand.Intn(m.Height + 1)}
}
