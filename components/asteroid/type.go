package asteroid

import "yioz.io/asteroid-miner/defs"

type Asteroid struct {
	defs.ObjectImage
	Id           int
	Size         int
	DrawPoints   []defs.Point
	Vector       defs.Point
	MaterialType MaterialType
}

type MaterialType int

const (
	None MaterialType = 0
	Gold
)
