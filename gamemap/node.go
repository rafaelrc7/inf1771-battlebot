package gamemap

func GetAdjacentPositions(m *Map, pos Coord) (adjs []Coord) {
	adjs = make([]Coord, 4)

	switch pos.D {
	case NORTH:
		adjs = append(adjs, Coord{pos.X, pos.Y, WEST})
		adjs = append(adjs, Coord{pos.X, pos.Y, EAST})
		if pos.Y < m.Height-1 {
			adjs = append(adjs, Coord{pos.X, pos.Y + 1, pos.D})
		}
		if pos.Y > 0 {
			adjs = append(adjs, Coord{pos.X, pos.Y - 1, pos.D})
		}
	case EAST:
		adjs = append(adjs, Coord{pos.X, pos.Y, NORTH})
		adjs = append(adjs, Coord{pos.X, pos.Y, SOUTH})
		if pos.X < m.Width-1 {
			adjs = append(adjs, Coord{pos.X + 1, pos.Y, pos.D})
		}
		if pos.X > 0 {
			adjs = append(adjs, Coord{pos.X - 1, pos.Y, pos.D})
		}
	case SOUTH:
		adjs = append(adjs, Coord{pos.X, pos.Y, EAST})
		adjs = append(adjs, Coord{pos.X, pos.Y, WEST})
		if pos.Y > 0 {
			adjs = append(adjs, Coord{pos.X, pos.Y - 1, pos.D})
		}
		if pos.Y < m.Height-1 {
			adjs = append(adjs, Coord{pos.X, pos.Y + 1, pos.D})
		}
	case WEST:
		adjs = append(adjs, Coord{pos.X, pos.Y, SOUTH})
		adjs = append(adjs, Coord{pos.X, pos.Y, NORTH})
		if pos.X > 0 {
			adjs = append(adjs, Coord{pos.X - 1, pos.Y, pos.D})
		}
		if pos.X < m.Width-1 {
			adjs = append(adjs, Coord{pos.X + 1, pos.Y, pos.D})
		}
	}

	return adjs
}
