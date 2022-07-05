package ai

import "math/rand"

func GetDecision() int {
	return rand.Intn(7)
}
