package ai

import (
	"container/heap"
	"math"

	"github.com/rafaelrc7/inf1771-battlebot/gamemap"
)

type node struct {
	prev            *node
	f, g, h, d      float64
	coord           gamemap.Coord
	index           int
	priority        float64
	inheap, visited bool
}

func manhattan(orig, dest gamemap.Coord) (dist float64) {
	return math.Abs(float64(orig.X-dest.X)) + math.Abs(float64(orig.Y-dest.Y))
}

func Astar(orig, dest gamemap.Coord, m *gamemap.Map) (actions []int, numActions int) {
	var nh nodeHeap
	heap.Init(&nh)

	nodes := make([][][]node, m.Width)
	for i := range nodes {
		nodes[i] = make([][]node, m.Height)
		for j := range nodes[i] {
			nodes[i][j] = make([]node, 4)
		}
	}

	curr := &nodes[orig.X][orig.Y][orig.D]

	curr.coord = orig
	heap.Push(&nh, curr)
	curr.inheap = true

	for curr.coord.X != dest.X || curr.coord.Y != dest.Y {
		if nh.Len() == 0 {
			return []int{}, 0
		}

		curr = heap.Pop(&nh).(*node)
		curr.inheap = false
		curr.visited = true

		if curr.coord.X != dest.X || curr.coord.Y != dest.Y {
			for _, adj := range m.GetAdjacentPositions(curr.coord) {
				adj_n := &nodes[adj.X][adj.Y][adj.D]

				if c := &m.Cells[adj.X][adj.Y]; c.Status == gamemap.WALL ||
					c.Status == gamemap.DANGEROUS ||
					c.Status == gamemap.TELEPORT ||
					c.Status == gamemap.HOLE {

					c.Visited = true
					continue

				}

				if adj_n.visited {
					continue
				} else if !adj_n.inheap {
					adj_n.coord = adj
					adj_n.d = 1
					adj_n.h = manhattan(adj, dest)

					adj_n.prev = curr
					adj_n.g = adj_n.d + adj_n.prev.g
					adj_n.f = adj_n.h + adj_n.g

					adj_n.priority = adj_n.f

					heap.Push(&nh, adj_n)
					adj_n.inheap = true
				} else if !adj_n.visited {
					g := adj_n.d + curr.g
					f := adj_n.h + g
					if nh.TryUpdate(adj_n, f) {
						adj_n.prev = curr
						adj_n.g = g
						adj_n.f = f
					}
				}

			}
		}
	}

	return path2actions(curr)
}

func peek() {
}

func path2actions(n *node) (actions []int, numActions int) {
	for n.prev != nil {
		curr := n.coord
		prev := n.prev.coord
		if curr.X == prev.X && curr.Y == prev.Y {
			switch curr.D {
			case gamemap.NORTH:
				switch prev.D {
				case gamemap.EAST:
					actions = append(actions, TURN_LEFT)
				case gamemap.WEST:
					actions = append(actions, TURN_RIGHT)
				}
			case gamemap.EAST:
				switch prev.D {
				case gamemap.NORTH:
					actions = append(actions, TURN_RIGHT)
				case gamemap.SOUTH:
					actions = append(actions, TURN_LEFT)
				}
			case gamemap.SOUTH:
				switch prev.D {
				case gamemap.EAST:
					actions = append(actions, TURN_RIGHT)
				case gamemap.WEST:
					actions = append(actions, TURN_LEFT)
				}
			case gamemap.WEST:
				switch prev.D {
				case gamemap.NORTH:
					actions = append(actions, TURN_LEFT)
				case gamemap.SOUTH:
					actions = append(actions, TURN_RIGHT)
				}
			}
		} else if curr.X < prev.X { /* E to W */
			if prev.D == gamemap.EAST {
				actions = append(actions, BACKWARD)
			} else {
				actions = append(actions, FORWARD)
			}
		} else if curr.X > prev.X { /* W to E */
			if prev.D == gamemap.EAST {
				actions = append(actions, FORWARD)
			} else {
				actions = append(actions, BACKWARD)
			}
		} else if curr.Y < prev.Y { /* N to S */
			if prev.D == gamemap.NORTH {
				actions = append(actions, FORWARD)
			} else {
				actions = append(actions, BACKWARD)
			}
		} else if curr.Y > prev.Y { /* S to N */
			if prev.D == gamemap.NORTH {
				actions = append(actions, BACKWARD)
			} else {
				actions = append(actions, FORWARD)
			}
		}

		n = n.prev
		numActions++
	}
	return actions, numActions
}
