package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image"
	"image/color"
	"log"
	"math"
	"math/rand"
	"strings"
	"yioz.io/asteroid-miner/resources/fonts"
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

func degToSinCos(deg int) (float32, float32) {
	var rotateDeg = float64(deg % 360)
	var rotateRad = (rotateDeg * math.Pi) / 180
	var cos = float32(math.Cos(rotateRad))
	var sin = float32(math.Sin(rotateRad))
	return sin, cos
}

func toPoint(speed float32, direction int) Point {
	var sin, cos = degToSinCos(direction)

	var vx = speed * cos
	var vy = speed * -sin

	return Point{vx, vy}
}
func rotate(x, y, sin, cos float32) (float32, float32) {
	return cos*x + sin*y, -sin*x + cos*y
}

func drawPlayer(screen *ebiten.Image, player *ObjectImage, number int) {
	var path vector.Path

	var width = player.width
	var height = player.height

	path.MoveTo(height/2, 0)
	path.LineTo(-height/2, +width/2)
	path.LineTo(-height/2, -width/2)
	path.LineTo(height/2, 0)

	var x2 = height * 0.8
	var y2 = width * 0.8
	path.LineTo(x2/2-1.5, 0)
	path.LineTo(-x2/2-1.5, +y2/2)
	path.LineTo(-x2/2-1.5, -y2/2)
	path.LineTo(x2/2-1.5, 0)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	var sin, cos = degToSinCos(player.direction)

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].DstX, vs[i].DstY = rotate(vs[i].DstX, vs[i].DstY, sin, cos)
		vs[i].DstX += centerX + player.x
		vs[i].DstY += centerY + player.y
		var tone = (0xdd - float32(number*0x22)) / 0xff
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone
	}
	player.drawnPoints[0].x = vs[0].DstX
	player.drawnPoints[0].y = vs[0].DstY
	player.drawnPoints[1].x = vs[1].DstX
	player.drawnPoints[1].y = vs[1].DstY
	player.drawnPoints[2].x = vs[2].DstX
	player.drawnPoints[2].y = vs[2].DstY
	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func updatePlayer(player *Player) {
	var v = toPoint(player.acceleration, player.image.direction)
	player.vector.x += v.x
	player.vector.y += v.y

	var speed = float32(math.Sqrt(math.Pow(float64(player.vector.x), 2) + math.Pow(float64(player.vector.y), 2)))
	if speed > maxSpeed {
		player.vector.x = player.vector.x * (maxSpeed / speed)
		player.vector.y = player.vector.y * (maxSpeed / speed)
	} else if speed > 0 {
		player.vector.x = player.vector.x * (float32(math.Max(float64(speed-0.1), 0)) / speed)
		player.vector.y = player.vector.y * (float32(math.Max(float64(speed-0.1), 0)) / speed)
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

var bulletId = 0

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
func shootBullet(player *Player) Bullet {
	vec := toPoint(maxSpeed*1.2, player.image.direction)
	player.vector.x -= vec.x / 8
	player.vector.y -= vec.y / 8
	var speed = float32(math.Sqrt(math.Pow(float64(player.vector.x), 2) + math.Pow(float64(player.vector.y), 2)))
	if speed > maxSpeed {
		player.vector.x = player.vector.x * (maxSpeed / speed)
		player.vector.y = player.vector.y * (maxSpeed / speed)
	}
	bulletId++
	return Bullet{bulletId, player.image.x, player.image.y, vec, 0}
}
func updateBullet(bullet *Bullet) {
	bullet.x += bullet.vector.x
	bullet.y += bullet.vector.y
	bullet.time += 1
}

func drawAsteroid(screen *ebiten.Image, image *ObjectImage) {
	var path vector.Path

	var width = image.width
	var height = image.height

	path.MoveTo(-height, -width)
	path.LineTo(-height, +width)
	path.LineTo(+height, +width)
	path.LineTo(+height, -width)

	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	var sin, cos = degToSinCos(image.direction)

	var vs, is = path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].DstX, vs[i].DstY = rotate(vs[i].DstX, vs[i].DstY, sin, cos)
		vs[i].DstX += centerX + image.x
		vs[i].DstY += centerY + image.y

		var tone = 0xaa / float32(0xff)
		vs[i].ColorR = tone
		vs[i].ColorG = tone
		vs[i].ColorB = tone

		image.drawnPoints[i].x = vs[i].DstX
		image.drawnPoints[i].y = vs[i].DstY
	}
	screen.DrawTriangles(vs, is, emptySubImage, op)
}
func updateAsteroid(asteroid *Asteroid) {
	asteroid.image.x += asteroid.vector.x
	asteroid.image.y += asteroid.vector.y

	for asteroid.image.x < -centerX {
		asteroid.image.x += screenWidth
	}
	for asteroid.image.y < -centerY {
		asteroid.image.y += screenHeight
	}
	for asteroid.image.x > centerX {
		asteroid.image.x -= screenWidth
	}
	for asteroid.image.y > centerY {
		asteroid.image.y -= screenHeight
	}

	detectCollisionByBullet(asteroid.image)
}
func detectCollisionByBullet(image ObjectImage) bool {
	var highestIndex = 0
	for i, v := range image.drawnPoints {
		if image.drawnPoints[highestIndex].y > v.y {
			highestIndex = i
		}
	}

	var topPoint = image.drawnPoints[highestIndex]
	var leftPoint = image.drawnPoints[(highestIndex+1)%4]
	var bottomPoint = image.drawnPoints[(highestIndex+2)%4]
	var rightPoint = image.drawnPoints[(highestIndex+3)%4]

	for _, v := range bullets {
		var x = v.x + centerX
		var y = v.y + centerY
		if v.time >= 1000 {
			continue
		}
		if y < topPoint.y || y > bottomPoint.y {
			continue
		}
		if x < leftPoint.x || x > rightPoint.x {
			continue
		}
		if topPoint.y == leftPoint.y || topPoint.y == rightPoint.y {
			v.time = 1000
			return true
		}

		var leftOfTop = x <= topPoint.x
		var leftOfBottom = x <= bottomPoint.x
		var topOfLeft = y <= leftPoint.y
		var topOfRight = y <= rightPoint.y

		var bottomOfTop = false
		if leftOfTop && topOfLeft {
			var xRate = 1 - (x-leftPoint.x)/(topPoint.x-leftPoint.x)
			if (y - topPoint.y) > (leftPoint.y-topPoint.y)*xRate {
				bottomOfTop = true
			}
		} else if !leftOfTop && topOfRight {
			var xRate = (x - topPoint.x) / (rightPoint.x - topPoint.x)
			if (y - topPoint.y) > (rightPoint.y-topPoint.y)*xRate {
				bottomOfTop = true
			}
		}
		if bottomOfTop {
			if leftOfBottom {
				var xRate = (x - leftPoint.x) / (bottomPoint.x - leftPoint.x)
				if (y - leftPoint.y) < (bottomPoint.y-leftPoint.y)*xRate {
					v.time = 1000
					return true
				}
			} else {
				var xRate = 1 - (x-bottomPoint.x)/(rightPoint.x-bottomPoint.x)
				if (y - rightPoint.y) < (bottomPoint.y-rightPoint.y)*xRate {
					v.time = 1000
					return true
				}
			}
		}
	}
	return false
}

func drawBulletCountUi(screen *ebiten.Image, usedBulletCount int) {
	var str = strings.Repeat("■", maxBullet-usedBulletCount) + strings.Repeat("□", usedBulletCount)
	text.Draw(screen, str, fontFace, centerX-(24*maxBullet/2), screenHeight-20, color.White)
}

type Game struct {
	counter int
}

const maxSpeed = 8
const maxBullet = 10
const maxAsteroidSizeSum = 20
const bulletTime = 120
const asteroidSpeed = 3

func (g *Game) Update() error {
	g.counter++

	var asteroidSizeSum = 0
	for _, v := range asteroids {
		asteroidSizeSum += v.size
	}

	if (maxAsteroidSizeSum-asteroidSizeSum) > 4 && g.counter%100 == 0 {
		var x, y float32 = 0, 0
		var vec = toPoint(float32(rand.Int()%(asteroidSpeed-1)+1), rand.Int())

		const defaultSize = 4
		var minSize = defaultSize*5 + 1
		var width = float32(rand.Int()%5 + minSize)
		var height = float32(rand.Int()%5 + minSize)
		if rand.Int()&1 == 1 {
			x = float32(rand.Int()%screenWidth - centerX)
			y = float32((rand.Int()&1)*screenHeight - centerY)
		} else {
			x = float32((rand.Int()&1)*screenWidth - centerX)
			y = float32(rand.Int()%screenHeight - centerY)
		}
		var img = ObjectImage{x, y, width, height, rand.Int() % 360, []Point{{0, 0}, {0, 0}, {0, 0}, {0, 0}}} // rand.Int() % 360}
		var asteroid = Asteroid{img, defaultSize, vec, None}
		asteroids = append(asteroids, &asteroid)
	}

	if player == nil {
		player = &Player{ObjectImage{0, 0, 20, 30, 90, []Point{{0, 0}, {0, 0}, {0, 0}}}, 0, Point{0, 0}, []ObjectImage{}}
	} else if g.counter%4 == 0 {
		var length = len(player.afterImage)
		if length <= 2 {
			player.afterImage = append([]ObjectImage{player.image}, player.afterImage...)
		} else {
			player.afterImage[2] = player.afterImage[1]
			player.afterImage[1] = player.afterImage[0]
			player.afterImage[0] = player.image
		}
	}

	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		player.image.direction = (player.image.direction + 3 + 360) % 360
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		player.image.direction = (player.image.direction - 3 + 360) % 360
	}
	if inpututil.KeyPressDuration(ebiten.KeyArrowUp) > 0 {
		player.acceleration = +0.5
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowDown) > 0 {
		player.acceleration = -0.5
	} else {
		player.acceleration = 0
	}

	updatePlayer(player)

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if len(bullets) < maxBullet {
			b := shootBullet(player)
			bullets = append(bullets, &b)
		}
	}

	for _, v := range asteroids {
		updateAsteroid(v)
	}
	for _, v := range bullets {
		updateBullet(v)
	}

	for i, v := range bullets {
		if v.time > bulletTime {
			bullets = append(bullets[0:i], bullets[i+1:]...)
		}
	}

	return nil
}

type Point struct {
	x, y float32
}
type Player struct {
	image        ObjectImage
	acceleration float32
	vector       Point
	afterImage   []ObjectImage
}
type ObjectImage struct {
	x, y, width, height float32
	direction           int
	drawnPoints         []Point
}
type Bullet struct {
	id     int
	x, y   float32
	vector Point
	time   int
}
type Asteroid struct {
	image  ObjectImage
	size   int
	vector Point
	aType  AsteroidType
}

type AsteroidType int

const (
	None AsteroidType = 0
	Gold
)

var player *Player
var bullets []*Bullet
var asteroids []*Asteroid

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	if player != nil {
		var length = len(player.afterImage)
		for i := length - 1; i >= 0; i-- {
			drawPlayer(screen, &player.afterImage[i], i+1)
		}
		drawPlayer(screen, &player.image, 0)
	}
	for _, v := range bullets {
		if v.time < bulletTime {
			drawBullet(screen, v)
		}
	}
	for _, v := range asteroids {
		drawAsteroid(screen, &v.image)
	}
	drawBulletCountUi(screen, len(bullets))

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f\nrotate: %d", ebiten.ActualTPS(), ebiten.ActualFPS(), player.image.direction))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

var fontFace font.Face

func main() {
	g := &Game{counter: 0}

	tt, err := opentype.Parse(fonts.PixelMplus12Regular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	fontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Asteroid Miner")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
