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

func DetectCollisionByBullet(image ObjectImage, bulletMap map[int]Point) []int {
	var hitBulletId []int

	var highestIndex = 0
	for i, v := range image.DrawnPoints {
		if image.DrawnPoints[highestIndex].Y > v.Y {
			highestIndex = i
		}
	}

	var topPoint = image.DrawnPoints[highestIndex]
	var leftPoint = image.DrawnPoints[(highestIndex+1)%4]
	var bottomPoint = image.DrawnPoints[(highestIndex+2)%4]
	var rightPoint = image.DrawnPoints[(highestIndex+3)%4]

	for bulletId, v := range bulletMap {
		var x = v.X + CenterX
		var y = v.Y + CenterY
		if y < topPoint.Y || y > bottomPoint.Y {
			continue
		}
		if x < leftPoint.X || x > rightPoint.X {
			continue
		}
		if topPoint.Y == leftPoint.Y || topPoint.Y == rightPoint.Y {
			hitBulletId = append(hitBulletId, bulletId)
			continue
		}

		var leftOfTop = x <= topPoint.X
		var leftOfBottom = x <= bottomPoint.X
		var topOfLeft = y <= leftPoint.Y
		var topOfRight = y <= rightPoint.Y

		var bottomOfTop = false
		if leftOfTop && topOfLeft {
			var xRate = 1 - (x-leftPoint.X)/(topPoint.X-leftPoint.X)
			if (y - topPoint.Y) > (leftPoint.Y-topPoint.Y)*xRate {
				bottomOfTop = true
			}
		} else if !leftOfTop && topOfRight {
			var xRate = (x - topPoint.X) / (rightPoint.X - topPoint.X)
			if (y - topPoint.Y) > (rightPoint.Y-topPoint.Y)*xRate {
				bottomOfTop = true
			}
		}
		if bottomOfTop {
			if leftOfBottom {
				var xRate = (x - leftPoint.X) / (bottomPoint.X - leftPoint.X)
				if (y - leftPoint.Y) < (bottomPoint.Y-leftPoint.Y)*xRate {
					hitBulletId = append(hitBulletId, bulletId)
					continue
				}
			} else {
				var xRate = 1 - (x-bottomPoint.X)/(rightPoint.X-bottomPoint.X)
				if (y - rightPoint.Y) < (bottomPoint.Y-rightPoint.Y)*xRate {
					hitBulletId = append(hitBulletId, bulletId)
					continue
				}
			}
		}
	}
	return hitBulletId
}
