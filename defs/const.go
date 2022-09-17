package defs

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

const (
	ScreenWidth  = 640
	ScreenHeight = 480
	CenterX      = ScreenWidth / 2
	CenterY      = ScreenHeight / 2
)

var (
	EmptyImage = ebiten.NewImage(3, 3)

	// EmptySubImage is an internal sub image of EmptyImage.
	// Use emptySubImage at DrawTriangles instead of emptyImage in order to avoid bleeding edges.
	EmptySubImage = EmptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)
