package main

const (
	north_p = iota
	east_p  = iota
	south_p = iota
	west_p  = iota
)

type Node struct {
	north, east, south, west *Node
	coord                    Coord
}

type Loc struct {
	coord Coord
	dir   int
}
