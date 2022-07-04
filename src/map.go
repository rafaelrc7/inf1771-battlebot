package main

const danger_decay = 1
const danger_base = 50

const (
	north = iota
	east  = iota
	south = iota
	west  = iota
)

const (
	wall      = iota
	dangerous = iota
	unknown   = iota
	teleport  = iota
	hole      = iota
	empty     = iota
	gold      = iota
	powerup   = iota
)

const (
	breeze    = 1 << iota
	flash     = 1 << iota
	redlight  = 1 << iota
	bluelight = 1 << iota
	weaklight = 1 << iota
)

type Cell struct {
	status       int
	senses       uint
	danger_level float64
	visited      bool
}

type Coord struct {
	x, y int
}

type Map struct {
	cells         [][]Cell
	height, width int

	gold_cells    map[Coord]uint64
	powerup_cells map[Coord]uint64
	danger_cells  map[Coord]*Cell
}

func newMap(h, w int) Map {
	var m Map

	m.height = h
	m.width = w

	m.cells = make([][]Cell, w)
	for i := range m.cells {
		m.cells[i] = make([]Cell, h)
	}

	m.gold_cells = make(map[Coord]uint64)
	m.powerup_cells = make(map[Coord]uint64)
	m.danger_cells = make(map[Coord]*Cell)

	return m
}

func tick(m *Map) {
	for key, val := range m.danger_cells {
		val.danger_level--
		if val.danger_level <= 0 {
			val.danger_level = 0
			delete(m.danger_cells, key)
		}
	}
	for key, val := range m.gold_cells {
		if val > 0 {
			m.gold_cells[key]--
		}
	}
	for key, val := range m.powerup_cells {
		if val > 0 {
			m.powerup_cells[key]--
		}
	}
}

func getAdjacentCells(m *Map, c Coord) (adjs []Coord) {
	adjs = make([]Coord, 4)

	if c.x > 0 {
		adjs = append(adjs, Coord{c.x - 1, c.y})
	}
	if c.y > 0 {
		adjs = append(adjs, Coord{c.x, c.y - 1})
	}
	if c.x < m.width-1 {
		adjs = append(adjs, Coord{c.x + 1, c.y})
	}
	if c.y < m.height-1 {
		adjs = append(adjs, Coord{c.x, c.y + 1})
	}

	return adjs
}

func visitCell(m *Map, c Coord, senses uint) {
	if m.cells[c.x][c.y].visited {
		return
	}

	adjs := getAdjacentCells(m, c)
	m.cells[c.x][c.y].visited = true
	m.cells[c.x][c.y].senses |= senses

	for _, ac := range adjs {
		status := &m.cells[ac.x][ac.y].status
		if *status == unknown || *status == dangerous {
			if isPossibleHole(m, ac) || isPossibleTeleport(m, ac) {
				*status = dangerous
			} else {
				*status = unknown
			}
		}
	}

	if senses&redlight != 0 {
		m.cells[c.x][c.y].status = powerup
	} else if senses&bluelight != 0 {
		m.cells[c.x][c.y].status = gold
	} else {
		m.cells[c.x][c.y].status = empty
	}
}

func addDanger(m *Map, c Coord) {
	m.cells[c.x][c.y].danger_level = danger_base
}

func isPossibleHole(m *Map, c Coord) bool {
	adjs := getAdjacentCells(m, c)
	for _, ac := range adjs {
		if m.cells[ac.x][ac.y].visited && m.cells[ac.x][ac.y].senses&breeze == 0 {
			return false
		}
	}
	return true
}

func isPossibleTeleport(m *Map, c Coord) bool {
	adjs := getAdjacentCells(m, c)
	for _, ac := range adjs {
		if m.cells[ac.x][ac.y].visited && m.cells[ac.x][ac.y].senses&flash == 0 {
			return false
		}
	}
	return true
}
