package main

import "math/rand"

const (
	turn_right   = iota
	turn_left    = iota
	forward      = iota
	backward     = iota
	attack       = iota
	take_gold    = iota
	take_powerup = iota
)

func getDecision() int {
	return rand.Intn(7)
}
