package player

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
	"yioz.io/asteroid-miner/defs"
)

var Instance *Player

func New() {
	Instance = &Player{
		ObjectImage: defs.ObjectImage{
			RectSize:    defs.Point{X: 20, Y: 30},
			Direction:   90,
			DrawnPoints: []defs.Point{{0, 0}, {0, 0}, {0, 0}},
		},
		Vector:     defs.Point{},
		AfterImage: []defs.ObjectImage{},
	}
}

func (player *Player) draw(screen *ebiten.Image, afterImageNumber int) {
	var image defs.ObjectImage
	if afterImageNumber == 0 {
		image = player.ObjectImage
	} else {
		image = player.AfterImage[afterImageNumber-1]
	}

	var path vector.Path

	var width = image.RectSize.X
	var height = image.RectSize.Y

	path.MoveTo(height/2, 0)
	path.LineTo(-height/2, +width/2)
	path.LineTo(-height/2, -width/2)

	var x2 = height * 0.8
	var y2 = width * 0.8
	path.MoveTo(x2/2-1.5, 0)
	path.LineTo(-x2/2-1.5, +y2/2)
	path.LineTo(-x2/2-1.5, -y2/2)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	var sin, cos = defs.DegToSinCos(image.Direction)

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].DstX, vs[i].DstY = defs.Rotate(vs[i].DstX, vs[i].DstY, sin, cos)
		vs[i].DstX += defs.CenterX + image.Position.X
		vs[i].DstY += defs.CenterY + image.Position.Y

		var tone = (0xdd - float32(afterImageNumber*0x22)) / 0xff
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone

		if i < 3 {
			image.DrawnPoints[i].X = vs[i].DstX
			image.DrawnPoints[i].Y = vs[i].DstY
		}
	}
	screen.DrawTriangles(vs, is, defs.EmptySubImage, op)
}

func (player *Player) Draw(screen *ebiten.Image) {
	var length = len(player.AfterImage)
	for i := length - 1; i >= 0; i-- {
		player.draw(screen, i+1)
	}
	player.draw(screen, 0)
}

func (player *Player) Update() {
	var v = defs.ToPoint(player.Acceleration, player.Direction)
	player.Vector.X += v.X
	player.Vector.Y += v.Y

	var speed = float32(math.Sqrt(math.Pow(float64(player.Vector.X), 2) + math.Pow(float64(player.Vector.Y), 2)))
	if speed > MaxSpeed {
		player.Vector.X = player.Vector.X * (MaxSpeed / speed)
		player.Vector.Y = player.Vector.Y * (MaxSpeed / speed)
	} else if speed > 0 {
		player.Vector.X = player.Vector.X * (float32(math.Max(float64(speed-0.1), 0)) / speed)
		player.Vector.Y = player.Vector.Y * (float32(math.Max(float64(speed-0.1), 0)) / speed)
	}

	player.Position.X += player.Vector.X
	player.Position.Y += player.Vector.Y

	for player.Position.X < -defs.CenterX {
		player.Position.X += defs.ScreenWidth
	}
	for player.Position.Y < -defs.CenterY {
		player.Position.Y += defs.ScreenHeight
	}
	for player.Position.X > defs.CenterX {
		player.Position.X -= defs.ScreenWidth
	}
	for player.Position.Y > defs.CenterY {
		player.Position.Y -= defs.ScreenHeight
	}
}

func (player *Player) UpdateAfterImage() {
	var length = len(player.AfterImage)
	if length <= 2 {
		player.AfterImage = append([]defs.ObjectImage{player.ObjectImage}, player.AfterImage...)
	} else {
		player.AfterImage[2] = player.AfterImage[1]
		player.AfterImage[1] = player.AfterImage[0]
		player.AfterImage[0] = player.ObjectImage
	}
}
