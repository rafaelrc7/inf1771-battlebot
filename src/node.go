package main

type NodeCoord struct {
	coord Coord
	dir   int
}

func getAdj(m *Map, pos NodeCoord) (adjs []NodeCoord) {
	adjs = make([]NodeCoord, 4)

	switch pos.dir {
	case north:
		adjs = append(adjs, NodeCoord{pos.coord, west})
		adjs = append(adjs, NodeCoord{pos.coord, east})
		if pos.coord.y < m.height-1 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x, pos.coord.y + 1}, pos.dir})
		}
		if pos.coord.y > 0 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x, pos.coord.y - 1}, pos.dir})
		}
	case east:
		adjs = append(adjs, NodeCoord{pos.coord, north})
		adjs = append(adjs, NodeCoord{pos.coord, south})
		if pos.coord.x < m.width-1 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x + 1, pos.coord.y}, pos.dir})
		}
		if pos.coord.x > 0 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x - 1, pos.coord.y}, pos.dir})
		}
	case south:
		adjs = append(adjs, NodeCoord{pos.coord, east})
		adjs = append(adjs, NodeCoord{pos.coord, west})
		if pos.coord.y > 0 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x, pos.coord.y - 1}, pos.dir})
		}
		if pos.coord.y < m.height-1 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x, pos.coord.y + 1}, pos.dir})
		}
	case west:
		adjs = append(adjs, NodeCoord{pos.coord, south})
		adjs = append(adjs, NodeCoord{pos.coord, north})
		if pos.coord.x > 0 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x - 1, pos.coord.y}, pos.dir})
		}
		if pos.coord.x < m.width-1 {
			adjs = append(adjs, NodeCoord{Coord{pos.coord.x + 1, pos.coord.y}, pos.dir})
		}
	}

	return adjs
}
