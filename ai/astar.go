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
	inheap, visited bool
}

func manhattan(orig, dest gamemap.Coord) (dist float64) {
	return math.Abs(float64(orig.X-dest.X)) + math.Abs(float64(orig.Y-dest.Y))
}

func Astar(orig, dest gamemap.Coord, m *gamemap.Map) (actions []int, numActions int) {
	actions = []int{}
	curr := gamemap.Coord{X: orig.X, Y: orig.Y, D: orig.D}
	queue := CoordHeapNew()

	nodes := make([][][]node, m.Width)
	for i := range nodes {
		nodes[i] = make([][]node, m.Height)
		for j := range nodes[i] {
			nodes[i][j] = make([]node, 4)
		}
	}

	nodes[curr.X][curr.Y][curr.D].coord = gamemap.Coord{X: curr.X, Y: curr.Y, D: curr.D}

	queue.PushCoord(curr, 0)
	nodes[curr.X][curr.Y][curr.D].inheap = true

	for curr.X != dest.X || curr.Y != dest.Y {
		if queue.Len() == 0 { /* no path */
			return []int{}, 0
		}
		curr = heap.Pop(queue).(gamemap.Coord)
		nodes[curr.X][curr.Y][curr.D].visited = true
		nodes[curr.X][curr.Y][curr.D].inheap = false

		if curr.X != dest.X || curr.Y != dest.Y {
			for _, adj := range m.GetAdjacentPositions(curr) {
				if m.Cells[adj.X][adj.Y].Status != gamemap.WALL &&
					m.Cells[adj.X][adj.Y].Status != gamemap.DANGEROUS &&
					m.Cells[adj.X][adj.Y].Status != gamemap.HOLE &&
					m.Cells[adj.X][adj.Y].Status != gamemap.TELEPORT {
					peek(adj, curr, dest, queue, &nodes, m)
				}
			}
		}
	}

	return path2actions(&nodes[curr.X][curr.Y][curr.D])
}

func peek(adj, curr, target gamemap.Coord, queue *CoordHeap, nodes *[][][]node, m *gamemap.Map) {
	node := &(*nodes)[adj.X][adj.Y][adj.D]
	prev := &(*nodes)[curr.X][curr.Y][curr.D]

	if node.visited {
		return
	} else if !node.inheap {
		node.coord = adj

		switch m.Cells[adj.X][adj.Y].Status {
		case gamemap.WALL:
		case gamemap.DANGEROUS:
		case gamemap.HOLE:
		case gamemap.TELEPORT:
			node.visited = true
			return

		default:
			node.d = 1
		}

		if m.Cells[adj.X][adj.Y].DangerLevel > 0 {
			node.d = 100
		}

		node.h = manhattan(target, node.coord)
		node.prev = prev
		node.g = node.d + node.prev.g
		node.f = node.h + node.g
		queue.PushCoord(adj, node.f)
	} else if !node.visited {
		g := node.d + prev.g
		f := node.h + g
		if queue.TryUpdate(node.coord, f) {
			node.prev = prev
			node.g = g
			node.f = f
		}
	}
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
