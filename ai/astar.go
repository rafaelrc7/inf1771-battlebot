package ai

import "math"
import "github.com/rafaelrc7/inf1771-battlebot/gamemap"

func manhattan(orig, dest gamemap.Coord) (dist float64) {
	return math.Abs(float64(orig.X-dest.X)) + math.Abs(float64(orig.Y-dest.Y))
}

func Astar(orig, dest gamemap.Coord) (actions []int) {
	return []int{}
}
