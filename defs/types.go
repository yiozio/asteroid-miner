package defs

type Point struct {
	X, Y float32
}

type BulletImage struct {
	Position Point
	Deleted  bool
}

type ObjectImage struct {
	Position    Point
	RectSize    Point
	Direction   int
	DrawnPoints []Point
}
