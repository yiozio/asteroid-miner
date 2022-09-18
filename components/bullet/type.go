package bullet

import "yioz.io/asteroid-miner/defs"

type Bullet struct {
	Id       int
	Position defs.Point
	Vector   defs.Point
	Time     int
}
