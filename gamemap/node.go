package gamemap

func (m *Map) GetAdjacentPositions(pos Coord) (adjs []Coord) {
	adjs = []Coord{}

	switch pos.D {
	case NORTH:
		if pos.Y > 0 {
			adjs = append(adjs, Coord{pos.X, pos.Y - 1, pos.D})
		}
		if pos.Y < m.Height-1 {
			adjs = append(adjs, Coord{pos.X, pos.Y + 1, pos.D})
		}
		adjs = append(adjs, Coord{pos.X, pos.Y, WEST})
		adjs = append(adjs, Coord{pos.X, pos.Y, EAST})
	case EAST:
		if pos.X < m.Width-1 {
			adjs = append(adjs, Coord{pos.X + 1, pos.Y, pos.D})
		}
		if pos.X > 0 {
			adjs = append(adjs, Coord{pos.X - 1, pos.Y, pos.D})
		}
		adjs = append(adjs, Coord{pos.X, pos.Y, NORTH})
		adjs = append(adjs, Coord{pos.X, pos.Y, SOUTH})
	case SOUTH:
		if pos.Y < m.Height-1 {
			adjs = append(adjs, Coord{pos.X, pos.Y + 1, pos.D})
		}
		if pos.Y > 0 {
			adjs = append(adjs, Coord{pos.X, pos.Y - 1, pos.D})
		}
		adjs = append(adjs, Coord{pos.X, pos.Y, EAST})
		adjs = append(adjs, Coord{pos.X, pos.Y, WEST})
	case WEST:
		if pos.X > 0 {
			adjs = append(adjs, Coord{pos.X - 1, pos.Y, pos.D})
		}
		if pos.X < m.Width-1 {
			adjs = append(adjs, Coord{pos.X + 1, pos.Y, pos.D})
		}
		adjs = append(adjs, Coord{pos.X, pos.Y, SOUTH})
		adjs = append(adjs, Coord{pos.X, pos.Y, NORTH})
	}

	return adjs
}

func (m *Map) GetForwardPosition(pos Coord) (f Coord) {
	switch pos.D {
	case NORTH:
		f = Coord{pos.X, pos.Y - 1, pos.D}
	case EAST:
		f = Coord{pos.X + 1, pos.Y, pos.D}
	case SOUTH:
		f = Coord{pos.X, pos.Y + 1, pos.D}
	case WEST:
		f = Coord{pos.X - 1, pos.Y, pos.D}
	}

	return f
}

func (m *Map) GetBackwardPosition(pos Coord) (f Coord) {
	switch pos.D {
	case SOUTH:
		f = Coord{pos.X, pos.Y - 1, pos.D}
	case WEST:
		f = Coord{pos.X + 1, pos.Y, pos.D}
	case NORTH:
		f = Coord{pos.X, pos.Y + 1, pos.D}
	case EAST:
		f = Coord{pos.X - 1, pos.Y, pos.D}
	}

	return f
}
