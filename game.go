package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/color"
	"strings"
	"yioz.io/asteroid-miner/components/asteroid"
	"yioz.io/asteroid-miner/components/bullet"
	"yioz.io/asteroid-miner/components/player"
	"yioz.io/asteroid-miner/defs"
)

func drawBulletCountUi(screen *ebiten.Image, usedBulletCount int) {
	var str = strings.Repeat("■", bullet.MaxBullet-usedBulletCount) + strings.Repeat("□", usedBulletCount)
	text.Draw(
		screen,
		str,
		fontFace,
		defs.CenterX-(24*bullet.MaxBullet/2),
		defs.ScreenHeight-20,
		color.White,
	)
}

type Game struct {
	counter int
}

func (g *Game) Update() error {
	g.counter++

	var asteroidSizeSum = 0
	for _, v := range asteroid.InstanceMap {
		asteroidSizeSum += v.Size
	}

	if (asteroid.MaxAsteroidSizeSum-asteroidSizeSum) > 4 && g.counter%100 == 0 {
		asteroid.Add()
	}

	if player.Instance == nil {
		player.New()
	} else if g.counter%4 == 0 {
		player.Instance.UpdateAfterImage()
	}

	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		player.Instance.Direction = (player.Instance.Direction + 3 + 360) % 360
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		player.Instance.Direction = (player.Instance.Direction - 3 + 360) % 360
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowUp) > 0 {
		player.Instance.Acceleration = +0.5
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowDown) > 0 {
		player.Instance.Acceleration = -0.5
	} else {
		player.Instance.Acceleration = 0
	}

	player.Instance.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if len(bullet.InstanceMap) < bullet.MaxBullet {
			bullet.Add(player.Instance)
		}
	}

	for aid, v := range asteroid.InstanceMap {
		v.Update()
		asteroid.InstanceMap[aid] = v
		var pointMap = map[int]defs.Point{}
		for bId, v := range bullet.InstanceMap {
			pointMap[bId] = v.Position
		}

		for i, v := range player.Instance.DrawnPoints {
			pointMap[-(i + 1)] = defs.Point{X: v.X - defs.CenterX, Y: v.Y - defs.CenterY}
		}

		for i, pointId := range defs.DetectCollisionByPoint(v.ObjectImage, pointMap) {
			if pointId < 0 {
				player.New()
				continue
			}
			if i == 0 {
				delete(asteroid.InstanceMap, aid)
			}
			bullet.Hit(pointId)
		}
	}

	for bId, v := range bullet.InstanceMap {
		v.Update()

		if v.Time > bullet.TimeToLive {
			delete(bullet.InstanceMap, bId)
		} else {
			bullet.InstanceMap[bId] = v
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if player.Instance != nil {
		player.Instance.Draw(screen)
	}
	for _, v := range bullet.InstanceMap {
		v.Draw(screen)
	}
	for _, v := range asteroid.InstanceMap {
		v.Draw(screen)
	}
	bullet.DrawHitEffect(screen)
	drawBulletCountUi(screen, len(bullet.InstanceMap))

	ebitenutil.DebugPrint(
		screen,
		fmt.Sprintf(
			"TPS: %0.2f\nFPS: %0.2f\nrotate: %d",
			ebiten.ActualTPS(),
			ebiten.ActualFPS(),
			player.Instance.Direction,
		),
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return defs.ScreenWidth, defs.ScreenHeight
}
