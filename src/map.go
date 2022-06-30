package main

const danger_decay = 1
const danger_base = 50

const (
	may_teleport = iota
	may_hole     = iota
	unknown      = iota
	empty        = iota
	teleport     = iota
	hole         = iota
	gold         = iota
	powerup      = iota
	wall         = iota
)

type Cell struct {
	status       int
	danger_level float64
	is_adj_safe  bool
	visited      bool
}

type Coord struct {
	x, y int64
}

type Map struct {
	cells [][]Cell

	gold_cells    map[Coord]struct{ last_collected uint64 }
	powerup_cells map[Coord]struct{ last_collected uint64 }
	danger_cells  map[Coord]*Cell
}

func newMap(h, w uint) Map {
	var m Map

	m.cells = make([][]Cell, h)
	for i := range m.cells {
		m.cells[i] = make([]Cell, w)
	}

	m.gold_cells = make(map[Coord]struct{ last_collected uint64 })
	m.powerup_cells = make(map[Coord]struct{ last_collected uint64 })
	m.danger_cells = make(map[Coord]*Cell)

	return m
}

func tick(m *Map) {
	for key, val := range m.danger_cells {
		(*val).danger_level--
		if (*val).danger_level <= 0 {
			(*val).danger_level = 0
			delete(m.danger_cells, key)
		}
	}
}

func set_may_hole(m *Map, c Coord) {
	if m.cells[c.x][c.y].is_adj_safe || m.cells[c.x][c.y].status > unknown {
		return
	}
	m.cells[c.x][c.y].status = may_hole
}

func set_may_teleport(m *Map, c Coord) {
	if m.cells[c.x][c.y].is_adj_safe || m.cells[c.x][c.y].status > unknown {
		return
	}
	m.cells[c.x][c.y].status = may_teleport
}
