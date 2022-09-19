package asteroid

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math/rand"
	"yioz.io/asteroid-miner/defs"
)

var asteroidId = 0

var InstanceMap = map[int]Asteroid{}

func Add() {
	var x, y float32 = 0, 0
	var vec = defs.ToPoint(float32(rand.Int()%(AsteroidSpeed-1)+1), rand.Int())

	const defaultSize = 4
	var minSize = defaultSize*5 + 1
	var width = float32(rand.Int()%5 + minSize)
	var height = float32(rand.Int()%5 + minSize)
	if rand.Int()&1 == 1 {
		x = float32(rand.Int()%defs.ScreenWidth - defs.CenterX)
		y = float32((rand.Int()&1)*defs.ScreenHeight - defs.CenterY)
	} else {
		x = float32((rand.Int()&1)*defs.ScreenWidth - defs.CenterX)
		y = float32(rand.Int()%defs.ScreenHeight - defs.CenterY)
	}
	var img = defs.ObjectImage{Position: defs.Point{X: x, Y: y}, RectSize: defs.Point{X: width, Y: height}, Direction: rand.Int() % 360, DrawnPoints: []defs.Point{{0, 0}, {0, 0}, {0, 0}, {0, 0}}}
	asteroidId += 1

	var drawPoints = []defs.Point{
		{X: -width, Y: -height},
		{X: float32(rand.Int()%int(width) - int(width)/2), Y: -height - float32(rand.Int()%20-10)},
		{X: +width, Y: -height},
		{X: +width + float32(rand.Int()%20-10), Y: float32(rand.Int()%int(height) - int(height)/2)},
		{X: +width, Y: +height},
		{X: float32(rand.Int()%int(width) - int(width)/2), Y: +height + float32(rand.Int()%20-10)},
		{X: -width, Y: +height},
		{X: -width - float32(rand.Int()%20-10), Y: float32(rand.Int()%int(height) - int(height)/2)},
	}

	InstanceMap[asteroidId] = Asteroid{ObjectImage: img, Id: asteroidId, Size: defaultSize, DrawPoints: drawPoints, Vector: vec, MaterialType: None}
}

func (asteroid *Asteroid) Draw(screen *ebiten.Image) {
	var path vector.Path

	for i, v := range asteroid.DrawPoints {
		if i == 0 {
			path.MoveTo(v.Y, v.X)
		} else {
			path.LineTo(v.Y, v.X)
		}
	}

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	var sin, cos = defs.DegToSinCos(asteroid.Direction)

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)

	for i := range vs {
		vs[i].DstX, vs[i].DstY = defs.Rotate(vs[i].DstX, vs[i].DstY, sin, cos)
		vs[i].DstX += defs.CenterX + asteroid.Position.X
		vs[i].DstY += defs.CenterY + asteroid.Position.Y

		var tone = 0xaa / float32(0xff)
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone

		if i%2 == 0 {
			asteroid.DrawnPoints[i/2].X = vs[i].DstX
			asteroid.DrawnPoints[i/2].Y = vs[i].DstY
		}
	}
	screen.DrawTriangles(vs, is, defs.EmptySubImage, op)
}

func (asteroid *Asteroid) Update() {
	asteroid.Position.X += asteroid.Vector.X
	asteroid.Position.Y += asteroid.Vector.Y

	for asteroid.Position.X < -defs.CenterX {
		asteroid.Position.X += defs.ScreenWidth
	}
	for asteroid.Position.Y < -defs.CenterY {
		asteroid.Position.Y += defs.ScreenHeight
	}
	for asteroid.Position.X > defs.CenterX {
		asteroid.Position.X -= defs.ScreenWidth
	}
	for asteroid.Position.Y > defs.CenterY {
		asteroid.Position.Y -= defs.ScreenHeight
	}
}
