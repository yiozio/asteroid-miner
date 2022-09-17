package player

import "yioz.io/asteroid-miner/defs"

type Player struct {
	defs.ObjectImage
	Acceleration float32
	Vector       defs.Point
	AfterImage   []defs.ObjectImage
}
