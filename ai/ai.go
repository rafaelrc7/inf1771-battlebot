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
	FLEEING
)

const (
	nothing = iota
	gold
	powerup
)

const minEnergy = 51
const maxTicksRunning = 15
const respawnTime = 150

type AI struct {
	State         int
	ActionStack   []int
	Coord         gamemap.Coord
	Dest          *gamemap.Coord
	Energy        int
	TimeRunnning  int
	TimeShooting  int
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

func (ai *AI) Think(mapChanged bool) {
	if ai.State == FETCHING_PU && ai.Energy == 100 {
		ai.State = EXPLORING
	}

	if ai.Energy < minEnergy {
		if ai.State == FETCHING_PU {
			return
		}
		dest := ai.findPUToFetch()

		if dest != nil {
			ai.State = FETCHING_PU
			ai.ActionStack = []int{}
			ai.Dest = dest
		}
		return
	}

	if ai.State == FLEEING && ai.TimeRunnning > maxTicksRunning {
		ai.State = EXPLORING
	}

	if ai.Observations&(gamemap.STEPS|gamemap.HIT) != 0 && !ai.EnemyDetected {
		if ai.State == FLEEING {
			return
		}
		ai.State = FLEEING
		ai.TimeRunnning = 0
		ai.ActionStack = []int{}
		ai.Dest = nil
		ai.Gamemap.AddDanger(ai.Coord)
		return
	}

	if ai.State != ATTACKING && ai.EnemyDetected {
		ai.State = ATTACKING
		ai.ActionStack = []int{}
		ai.Dest = nil
		ai.TimeShooting = 0
		return
	}

	if gold := ai.findGoldToFetch(); gold != nil {
		if ai.State == FETCHING_GOLD {
			return
		}
		ai.State = FETCHING_GOLD
		ai.ActionStack = []int{}
		ai.Dest = &gamemap.Coord{X: gold.X, Y: gold.Y}
		return
	}

	switch ai.State {
	case ATTACKING:
		if !ai.EnemyDetected {
			ai.State = EXPLORING
			return
		}
	default:
		ai.State = EXPLORING
	}
}

func (ai *AI) GetDecision(mapChanged bool) int {
	if mapChanged {
		ai.ActionStack = []int{}
	}

	if ai.TimeShooting == 10 {
		return TURN_LEFT
	}

	switch ai.State {
	case STOP:
		return NOTHING

	case ATTACKING:
		return ATTACK

	case EXPLORING:
		if ai.Observations&(gamemap.BLUELIGHT) != 0 {
			*ai.Gamemap.GoldCells[gamemap.CellC{X: ai.Coord.X, Y: ai.Coord.Y}] = respawnTime
			return TAKE
		}
		if ai.Observations&(gamemap.REDLIGHT) != 0 && ai.Energy < 100 {
			*ai.Gamemap.PowerupCells[gamemap.CellC{X: ai.Coord.X, Y: ai.Coord.Y}] = respawnTime
			return TAKE
		}
		if len(ai.ActionStack) == 0 {
			if ai.Dest == nil || ai.Dest.X == ai.Coord.X && ai.Dest.Y == ai.Coord.Y {
				dest := FindUnexplored(ai.Gamemap, ai.Coord)
				ai.Dest = &dest
			}
			ai.ActionStack, _ = Astar(ai.Coord, *ai.Dest, ai.Gamemap)
		}
		if len(ai.ActionStack) == 0 {
			ai.Dest = nil
			return NOTHING
		}

		l := len(ai.ActionStack)
		action := ai.ActionStack[l-1]
		ai.ActionStack = ai.ActionStack[:l-1]
		return action

	case FETCHING_GOLD:
		if r := ai.take(); r != nothing {
			if r == gold {
				ai.State = STOP
			}
			return TAKE
		}

		if len(ai.ActionStack) == 0 {
			if ai.Dest == nil {
				return NOTHING
			}
			ai.ActionStack, _ = Astar(ai.Coord, *ai.Dest, ai.Gamemap)
		}
		if len(ai.ActionStack) == 0 {
			ai.State = EXPLORING
			return NOTHING
		}

		l := len(ai.ActionStack)
		action := ai.ActionStack[l-1]
		ai.ActionStack = ai.ActionStack[:l-1]
		return action

	case FETCHING_PU:
		if r := ai.take(); r != nothing {
			if r == powerup {
				ai.State = STOP
			}
			return TAKE
		}

		if len(ai.ActionStack) == 0 {
			if ai.Dest == nil {
				ai.Dest = ai.findPUToFetch()
			}
			if ai.Dest == nil {
				ai.State = EXPLORING
			}
			ai.ActionStack, _ = Astar(ai.Coord, *ai.Dest, ai.Gamemap)
		}
		if len(ai.ActionStack) == 0 {
			ai.State = EXPLORING
			return NOTHING
		}

		l := len(ai.ActionStack)
		action := ai.ActionStack[l-1]
		ai.ActionStack = ai.ActionStack[:l-1]
		return action

	case FLEEING:
		ai.TimeRunnning++
		if len(ai.ActionStack) == 0 {
			if ai.Dest == nil {
				ai.Dest = ai.findPUToFetch()
			}
			if ai.Dest == nil {
				ai.Dest = &gamemap.Coord{X: rand.Intn(ai.Gamemap.Width + 1), Y: rand.Intn(ai.Gamemap.Height + 1)}
			}
			ai.ActionStack, _ = Astar(ai.Coord, *ai.Dest, ai.Gamemap)
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

func (ai *AI) take() int {
	c := gamemap.CellC{X: ai.Coord.X, Y: ai.Coord.Y}
	if ai.Observations&(gamemap.BLUELIGHT) != 0 {
		if ai.Gamemap.GoldCells[c] == nil {
			ai.Gamemap.GoldCells[c] = new(int)
		}
		*ai.Gamemap.GoldCells[c] = respawnTime
		return gold
	} else if ai.Observations&(gamemap.REDLIGHT) != 0 {
		if ai.Gamemap.PowerupCells[c] == nil {
			ai.Gamemap.PowerupCells[c] = new(int)
		}
		if ai.Energy < 100 {
			*ai.Gamemap.PowerupCells[c] = respawnTime
			return powerup
		}
	}

	return nothing
}

func (ai *AI) findGoldToFetch() *gamemap.Coord {
	cells := make(chan *gamemap.Coord, 10)

	workers := 0
	for k, v := range ai.Gamemap.GoldCells {
		if v != nil {
			go func(cell gamemap.Coord, spawntime int) {
				_, dist := Astar(ai.Coord, cell, ai.Gamemap)
				fmt.Printf("GOLD: (%d, %d) %d/%d\n", cell.X, cell.Y, dist, spawntime)
				if dist > 0 && dist >= spawntime {
					cells <- &cell
				} else {
					cells <- nil
				}
			}(gamemap.Coord{X: k.X, Y: k.Y}, *v)
			workers++
		}
	}

	for ; workers > 0; workers-- {
		cell := <-cells
		if cell != nil {
			return cell
		}
	}

	return nil
}

func (ai *AI) findPUToFetch() *gamemap.Coord {
	cells := make(chan *gamemap.Coord, 10)

	workers := 0
	for k, v := range ai.Gamemap.PowerupCells {
		if v != nil {
			go func(cell gamemap.Coord, spawntime int) {
				_, dist := Astar(ai.Coord, cell, ai.Gamemap)
				if dist > 0 && dist >= spawntime {
					cells <- &cell
				} else {
					cells <- nil
				}
			}(gamemap.Coord{X: k.X, Y: k.Y}, *v)
			workers++
		}
	}

	for ; workers > 0; workers-- {
		cell := <-cells
		if cell != nil {
			return cell
		}
	}

	return nil
}

func FindUnexplored(m *gamemap.Map, c gamemap.Coord) gamemap.Coord {
	for len(m.ExploreStack) > 0 {
		fmt.Println(m.ExploreStack)
		if c, s := m.StackPop(); s {
			return c
		}
	}

	return gamemap.Coord{X: rand.Intn(m.Width + 1), Y: rand.Intn(m.Height + 1)}
}
