package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"yioz.io/asteroid-miner/defs"
	"yioz.io/asteroid-miner/resources/fonts"
)

func init() {
	defs.EmptyImage.Fill(color.White)
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

	ebiten.SetWindowSize(defs.ScreenWidth, defs.ScreenHeight)
	ebiten.SetWindowTitle("Asteroid Miner")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
