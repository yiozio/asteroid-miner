package defs

import (
	"math"
)

func DegToSinCos(deg int) (float32, float32) {
	var rotateDeg = float64(deg % 360)
	var rotateRad = (rotateDeg * math.Pi) / 180
	var cos = float32(math.Cos(rotateRad))
	var sin = float32(math.Sin(rotateRad))
	return sin, cos
}

func ToPoint(speed float32, direction int) Point {
	var sin, cos = DegToSinCos(direction)

	var vx = speed * cos
	var vy = speed * -sin

	return Point{vx, vy}
}

func Rotate(x, y, sin, cos float32) (float32, float32) {
	return cos*x + sin*y, -sin*x + cos*y
}

func DetectCollisionByPoint(image ObjectImage, pointMap map[int]Point) []int {
	var hitBulletIds []int

	for bid, bullet := range pointMap {
		count := 0
		for i, point := range image.DrawnPoints {
			nextPoint := image.DrawnPoints[(i+1)%len(image.DrawnPoints)]

			y := bullet.Y + CenterY
			x := bullet.X + CenterX
			if point.Y < y && nextPoint.Y >= y || point.Y >= y && nextPoint.Y < y {
				rate := (y - point.Y) / (nextPoint.Y - point.Y)

				if x < point.X+(rate*(nextPoint.X-point.X)) {
					count += 1
				}
			}
		}
		if count%2 == 1 {
			hitBulletIds = append(hitBulletIds, bid)
		}
	}

	return hitBulletIds
}
