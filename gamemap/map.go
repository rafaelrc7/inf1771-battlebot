package gamemap

import "fmt"

const danger_decay = 1
const danger_base = 50

type Cell struct {
	Status      int
	Senses      uint
	DangerLevel float64
	Visited     bool
}

type Coord struct {
	X, Y, D int
}

type CellC struct {
	X, Y int
}

type Collectable struct {
	IsCollectable bool
	TicksToSpawn  int
}

type Map struct {
	Cells         [][]Cell
	Height, Width int

	GoldCells    map[CellC]*int
	PowerupCells map[CellC]*int
	DangerCells  map[CellC]*Cell
}

func NewMap(h, w int) *Map {
	var m Map

	m.Height = h
	m.Width = w

	m.Cells = make([][]Cell, w+1)
	for i := range m.Cells {
		m.Cells[i] = make([]Cell, h+1)
	}

	m.GoldCells = make(map[CellC]*int)
	m.PowerupCells = make(map[CellC]*int)
	m.DangerCells = make(map[CellC]*Cell)

	return &m
}

func (m *Map) Tick() bool {
	hasChanged := false
	for key, val := range m.DangerCells {
		val.DangerLevel--
		if val.DangerLevel <= 0 {
			val.DangerLevel = 0
			delete(m.DangerCells, key)
			return true
		}
	}
	for key, val := range m.GoldCells {
		if val != nil && *val > 0 {
			*m.GoldCells[key]--
		}
	}
	for key, val := range m.PowerupCells {
		if val != nil && *val > 0 {
			*m.PowerupCells[key]--
		}
	}
	return hasChanged
}

func (m *Map) GetAdjacentCells(c Coord) (adjs []Coord) {
	adjs = []Coord{}

	if c.X > 0 {
		adjs = append(adjs, Coord{c.X - 1, c.Y, 0})
	}
	if c.Y > 0 {
		adjs = append(adjs, Coord{c.X, c.Y - 1, 0})
	}
	if c.X < m.Width-1 {
		adjs = append(adjs, Coord{c.X + 1, c.Y, 0})
	}
	if c.Y < m.Height-1 {
		adjs = append(adjs, Coord{c.X, c.Y + 1, 0})
	}

	return adjs
}

func (m *Map) VisitCell(c Coord, senses uint) bool {
	hasChanged := false
	if m.Cells[c.X][c.Y].Visited {
		return hasChanged
	}

	adjs := m.GetAdjacentCells(c)
	m.Cells[c.X][c.Y].Visited = true
	m.Cells[c.X][c.Y].Senses |= senses

	for _, ac := range adjs {
		status := &m.Cells[ac.X][ac.Y].Status
		if *status == UNKNOWN || *status == DANGEROUS {
			if m.isPossibleHole(ac) || m.isPossibleTeleport(ac) {
				if *status != DANGEROUS {
					*status = DANGEROUS
					hasChanged = true
				}
			} else {
				if *status != SAFE {
					*status = SAFE
					hasChanged = true
				}
			}
		}
	}

	if senses&REDLIGHT != 0 {
		if m.Cells[c.X][c.Y].Status != POWERUP {
			m.Cells[c.X][c.Y].Status = POWERUP
			m.PowerupCells[CellC{X: c.X, Y: c.Y}] = new(int)
			hasChanged = true
		}
	} else if senses&BLUELIGHT != 0 {
		if m.Cells[c.X][c.Y].Status != GOLD {
			m.Cells[c.X][c.Y].Status = GOLD
			m.GoldCells[CellC{X: c.X, Y: c.Y}] = new(int)
			hasChanged = true
		}
	} else {
		if m.Cells[c.X][c.Y].Status != EMPTY {
			m.Cells[c.X][c.Y].Status = EMPTY
			hasChanged = true
		}
	}

	return hasChanged
}

func (m *Map) MarkWall(c Coord, forward bool) bool {
	hasChanged := false
	if forward {
		forward := m.GetForwardPosition(c)
		if forward.X < 0 || forward.Y < 0 {
			return hasChanged
		}
		if m.Cells[forward.X][forward.Y].Status != WALL {
			m.Cells[forward.X][forward.Y].Status = WALL
			hasChanged = true
		}
	} else {
		backward := m.GetBackwardPosition(c)
		if backward.X < 0 || backward.Y < 0 {
			return hasChanged
		}
		if m.Cells[backward.X][backward.Y].Status != WALL {
			m.Cells[backward.X][backward.Y].Status = WALL
			hasChanged = true
		}
	}

	return hasChanged
}

func (m *Map) AddDanger(c Coord) {
	m.Cells[c.X][c.Y].DangerLevel = danger_base
	adjs := m.GetAdjacentCells(c)
	for _, adj := range adjs {
		m.Cells[adj.X][adj.Y].DangerLevel = danger_base
	}
}

func (m *Map) Print(pos Coord) {
	fmt.Println("\nMAP:")
	for y := 0; y < m.Height; y++ {
		for x := 0; x < m.Width; x++ {
			if x == pos.X && y == pos.Y {
				switch pos.D {
				case NORTH:
					fmt.Print("^")
				case SOUTH:
					fmt.Print("v")
				case EAST:
					fmt.Print(">")
				case WEST:
					fmt.Print("<")
				}
			} else {
				switch m.Cells[x][y].Status {
				case WALL:
					fmt.Print("X")
				case DANGEROUS:
					fmt.Print("!")
				case UNKNOWN:
					fmt.Print("?")
				case SAFE:
					fmt.Print(",")
				case TELEPORT:
					fmt.Print("T")
				case HOLE:
					fmt.Print("O")
				case EMPTY:
					fmt.Print(".")
				case GOLD:
					fmt.Print("G")
				case POWERUP:
					fmt.Print("P")
				}
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("---------------------------------------------\n\n")
}

func (m *Map) isPossibleHole(c Coord) bool {
	adjs := m.GetAdjacentCells(c)
	for _, ac := range adjs {
		if m.Cells[ac.X][ac.Y].Visited && m.Cells[ac.X][ac.Y].Senses&BREEZE == 0 {
			return false
		}
	}
	return true
}

func (m *Map) isPossibleTeleport(c Coord) bool {
	adjs := m.GetAdjacentCells(c)
	for _, ac := range adjs {
		if m.Cells[ac.X][ac.Y].Visited && m.Cells[ac.X][ac.Y].Senses&FLASH == 0 {
			return false
		}
	}
	return true
}
