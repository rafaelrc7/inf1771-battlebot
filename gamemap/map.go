package gamemap

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

type Map struct {
	Cells         [][]Cell
	Height, Width int

	GoldCells    map[Coord]uint64
	PowerupCells map[Coord]uint64
	DangerCells  map[Coord]*Cell
}

func NewMap(h, w int) *Map {
	var m Map

	m.Height = h
	m.Width = w

	m.Cells = make([][]Cell, w)
	for i := range m.Cells {
		m.Cells[i] = make([]Cell, h)
	}

	m.GoldCells = make(map[Coord]uint64)
	m.PowerupCells = make(map[Coord]uint64)
	m.DangerCells = make(map[Coord]*Cell)

	return &m
}

func (m *Map) Tick() {
	for key, val := range m.DangerCells {
		val.DangerLevel--
		if val.DangerLevel <= 0 {
			val.DangerLevel = 0
			delete(m.DangerCells, key)
		}
	}
	for key, val := range m.GoldCells {
		if val > 0 {
			m.GoldCells[key]--
		}
	}
	for key, val := range m.PowerupCells {
		if val > 0 {
			m.PowerupCells[key]--
		}
	}
}

func (m *Map) GetAdjacentCells(c Coord) (adjs []Coord) {
	adjs = make([]Coord, 4)

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

func (m *Map) VisitCell(c Coord, senses uint) {
	if m.Cells[c.X][c.Y].Visited {
		return
	}

	adjs := m.GetAdjacentCells(c)
	m.Cells[c.X][c.Y].Visited = true
	m.Cells[c.X][c.Y].Senses |= senses

	for _, ac := range adjs {
		status := &m.Cells[ac.X][ac.Y].Status
		if *status == UNKNOWN || *status == DANGEROUS {
			if m.isPossibleHole(ac) || m.isPossibleTeleport(ac) {
				*status = DANGEROUS
			} else {
				*status = UNKNOWN
			}
		}
	}

	if senses&REDLIGHT != 0 {
		m.Cells[c.X][c.Y].Status = POWERUP
	} else if senses&BLUELIGHT != 0 {
		m.Cells[c.X][c.Y].Status = GOLD
	} else {
		m.Cells[c.X][c.Y].Status = EMPTY
	}
}

func (m *Map) AddDanger(c Coord) {
	m.Cells[c.X][c.Y].DangerLevel = danger_base
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