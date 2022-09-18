package defs

type Point struct {
	X, Y float32
}

type ObjectImage struct {
	Position    Point
	RectSize    Point
	Direction   int
	DrawnPoints []Point
}
