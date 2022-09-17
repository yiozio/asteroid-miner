package bullet

import "yioz.io/asteroid-miner/defs"

type Bullet struct {
	defs.BulletImage
	Id     int
	Vector defs.Point
	Time   int
}
