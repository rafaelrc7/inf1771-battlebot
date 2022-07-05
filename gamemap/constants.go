package gamemap

const (
	NORTH = iota
	EAST  = iota
	SOUTH = iota
	WEST  = iota
)

const (
	WALL      = iota
	DANGEROUS = iota
	UNKNOWN   = iota
	TELEPORT  = iota
	HOLE      = iota
	EMPTY     = iota
	GOLD      = iota
	POWERUP   = iota
)

const (
	BREEZE    = 1 << iota
	FLASH     = 1 << iota
	REDLIGHT  = 1 << iota
	BLUELIGHT = 1 << iota
	WEAKLIGHT = 1 << iota
)
