package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image"
	"image/color"
	"log"
	"math"
)

var (
	emptyImage = ebiten.NewImage(3, 3)

	// emptySubImage is an internal sub image of emptyImage.
	// Use emptySubImage at DrawTriangles instead of emptyImage in order to avoid bleeding edges.
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

const (
	screenWidth  = 640
	screenHeight = 480
	centerX      = screenWidth / 2
	centerY      = screenHeight / 2
)

func drawPlayer(screen *ebiten.Image, player *PlayerImage, number int) {
	var path vector.Path

	const width = 20
	const height = 30

	path.MoveTo(height/2, 0)
	path.LineTo(-height/2, +width/2)
	path.LineTo(-height/2, -width/2)
	path.LineTo(height/2, 0)

	const x2 = height * 0.8
	const y2 = width * 0.8
	path.LineTo(x2/2-1.5, 0)
	path.LineTo(-x2/2-1.5, +y2/2)
	path.LineTo(-x2/2-1.5, -y2/2)
	path.LineTo(x2/2-1.5, 0)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	var rotateDeg = float64(player.direction % 360)
	var rotateRad = (rotateDeg * math.Pi) / 180
	var cos = float32(math.Cos(rotateRad))
	var sin = float32(math.Sin(rotateRad))

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		var _x = vs[i].DstX
		var _y = vs[i].DstY
		vs[i].DstX = cos*_x + sin*_y
		vs[i].DstY = -sin*_x + cos*_y
		vs[i].DstX += centerX + player.x
		vs[i].DstY += centerY + player.y
		var tone = (0xdd - float32(number*0x22)) / 0xff
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func updatePlayer(player *Player) {
	var v = toVector(player.acceleration, player.image.direction)
	player.vector.x += v.x
	player.vector.y += v.y

	var speed = float32(math.Sqrt(math.Pow(float64(player.vector.x), 2) + math.Pow(float64(player.vector.y), 2)))
	if speed > maxSpeed {
		player.vector.x = player.vector.x * (maxSpeed / speed)
		player.vector.y = player.vector.y * (maxSpeed / speed)
	}

	player.image.x += player.vector.x
	player.image.y += player.vector.y

	for player.image.x < -centerX {
		player.image.x += screenWidth
	}
	for player.image.y < -centerY {
		player.image.y += screenHeight
	}
	for player.image.x > centerX {
		player.image.x -= screenWidth
	}
	for player.image.y > centerY {
		player.image.y -= screenHeight
	}
}

func drawBullet(screen *ebiten.Image, bullet *Bullet) {
	var path vector.Path

	const size = 2

	path.MoveTo(0, 0)
	path.Arc(0, 0, size, 0, 360*math.Pi/180, vector.Clockwise)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].DstX += centerX + bullet.x
		vs[i].DstY += centerY + bullet.y

		var tone = 0xaa / float32(0xff)
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}
func updateBullet(bullet *Bullet) {
	bullet.x += bullet.vector.x
	bullet.y += bullet.vector.y
}

type Game struct {
	counter int
}

const maxSpeed = 10

func (g *Game) Update() error {
	g.counter++

	if player == nil {
		player = &Player{PlayerImage{0, 0, 90}, 0, Vector{0, 0}, []PlayerImage{}}
	} else if g.counter%10 == 0 {
		var length = len(player.afterImage)
		if length <= 2 {
			player.afterImage = append([]PlayerImage{player.image}, player.afterImage...)
		} else {
			player.afterImage[2] = player.afterImage[1]
			player.afterImage[1] = player.afterImage[0]
			player.afterImage[0] = player.image
		}
	}

	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		player.image.direction = (player.image.direction + 2 + 360) % 360
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		player.image.direction = (player.image.direction - 2 + 360) % 360
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowUp) > 0 {
		player.acceleration = +0.5
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowDown) > 0 {
		player.acceleration = -0.5
	} else {
		player.acceleration = 0
	}

	updatePlayer(player)

	if bullet == nil {
		bullet = &Bullet{player.image.x, player.image.y, toVector(maxSpeed*1.2, player.image.direction)}
	}
	updateBullet(bullet)
	if bullet.x < -screenWidth/2 || bullet.x > screenWidth/2 || bullet.y < -screenHeight/2 || bullet.y > screenHeight/2 {
		bullet = nil
	}

	return nil
}

type Vector struct {
	x, y float32
}
type Player struct {
	image        PlayerImage
	acceleration float32
	vector       Vector
	afterImage   []PlayerImage
}
type PlayerImage struct {
	x, y      float32
	direction int
}
type Bullet struct {
	x, y   float32
	vector Vector
}

var player *Player
var bullet *Bullet

func toVector(speed float32, direction int) Vector {
	var deg = float64(direction % 360)
	var rad = deg * math.Pi / 180
	var cos = float32(math.Cos(rad))
	var sin = float32(math.Sin(rad))

	var vx = speed * cos
	var vy = speed * -sin

	return Vector{vx, vy}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if player != nil {
		var length = len(player.afterImage)
		for i := length - 1; i >= 0; i-- {
			drawPlayer(screen, &player.afterImage[i], i+1)
		}
		drawPlayer(screen, &player.image, 0)
	}
	if bullet != nil {
		drawBullet(screen, bullet)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nrotate: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), player.image.direction))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g := &Game{counter: 0}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Vector (Ebiten Demo)")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
