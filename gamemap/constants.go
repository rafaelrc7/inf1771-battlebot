package gamemap

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

const (
	UNKNOWN = iota
	SAFE
	WALL
	DANGEROUS
	TELEPORT
	HOLE
	EMPTY
	GOLD
	POWERUP
)

const (
	BREEZE    = 1 << iota
	FLASH     = 1 << iota
	REDLIGHT  = 1 << iota
	BLUELIGHT = 1 << iota
	WEAKLIGHT = 1 << iota
	STEPS     = 1 << iota
	BLOCKED   = 1 << iota
)
