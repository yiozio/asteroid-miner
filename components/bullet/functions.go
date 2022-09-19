package bullet

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"math"
	"yioz.io/asteroid-miner/components/player"
	"yioz.io/asteroid-miner/defs"
)

var bulletId = 0

var InstanceMap = map[int]Bullet{}
var hitEffect = map[int]Bullet{}

func Add(_player *player.Player) {
	vec := defs.ToPoint(player.MaxSpeed*1.2, _player.Direction)
	_player.Vector.X -= vec.X / 8
	_player.Vector.Y -= vec.Y / 8
	var speed = float32(math.Sqrt(math.Pow(float64(_player.Vector.X), 2) + math.Pow(float64(_player.Vector.Y), 2)))
	if speed > player.MaxSpeed {
		_player.Vector.X = _player.Vector.X * (player.MaxSpeed / speed)
		_player.Vector.Y = _player.Vector.Y * (player.MaxSpeed / speed)
	}
	bulletId++
	InstanceMap[bulletId] = Bullet{Position: _player.Position, Id: bulletId, Vector: vec}
}
func Hit(bulletId int) {
	var bullet = InstanceMap[bulletId]
	bullet.Time = 0
	hitEffect[bulletId] = bullet
	delete(InstanceMap, bulletId)
}
func DrawHitEffect(screen *ebiten.Image) {
	for id, v := range hitEffect {
		if v.Time > 6 {
			delete(hitEffect, id)
			continue
		}
		v.Time += 1
		hitEffect[id] = v

		var path vector.Path
		var size = float32(math.Sqrt(float64(v.Time+1) * 100))
		path.MoveTo(v.Position.X+defs.CenterX, v.Position.Y+defs.CenterY)
		path.Arc(v.Position.X+defs.CenterX, v.Position.Y+defs.CenterY, size, 0, 360*math.Pi/180, vector.Clockwise)
		path.MoveTo(v.Position.X+defs.CenterX, v.Position.Y+defs.CenterY)
		path.Arc(v.Position.X+defs.CenterX, v.Position.Y+defs.CenterY, size-2, 0, 360*math.Pi/180, vector.Clockwise)

		op := &ebiten.DrawTrianglesOptions{
			FillRule: ebiten.EvenOdd,
		}

		var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			var tone = 0xaa / float32(0xff)
			vs[i].ColorR = tone
			vs[i].ColorG = tone
			vs[i].ColorB = tone
		}
		screen.DrawTriangles(vs, is, defs.EmptySubImage, op)
	}
}
func (bullet *Bullet) Draw(screen *ebiten.Image) {
	var path vector.Path

	const size = 2

	path.MoveTo(0, 0)
	path.Arc(0, 0, size, 0, 360*math.Pi/180, vector.Clockwise)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].DstX += defs.CenterX + bullet.Position.X
		vs[i].DstY += defs.CenterY + bullet.Position.Y

		var tone = 0xaa / float32(0xff)
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone
	}
	screen.DrawTriangles(vs, is, defs.EmptySubImage, op)
}
func (bullet *Bullet) Update() {
	bullet.Position.X += bullet.Vector.X
	bullet.Position.Y += bullet.Vector.Y
	bullet.Time += 1
}
