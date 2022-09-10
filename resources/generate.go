//go:generate file2byteslice -package=fonts -input=./fonts/PixelMplus12-Regular.ttf -output=./fonts/PixelMplus12Regular.go -var=PixelMplus12Regular_ttf
//go:generate gofmt -s -w .

package resources

import (
	_ "github.com/hajimehoshi/file2byteslice"
)
