package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/coloring"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/explorer"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/fractal"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"
)

const (
	screenWidth  int     = 3840 // Adjust to your screen resolution
	screenHeight int     = 2160
	xMin         float64 = -3.0
	xMax         float64 = 1.2
	yMin         float64 = -1.2
	yMax         float64 = 1.2
)

func main() {
	// Initialize the pixel buffer
	screen := explorer.NewScreen(screenWidth, screenHeight)

	iterator := fractal.Iterator{
		Equation:  fractal.Mandlebrot,
		AbsLimit:  2.0,
		IterLimit: 64,
	}
	pallete := coloring.NewPallete(
		iterator.IterLimit,
		color.RGBA{R: 0, G: 0, B: 0},       // Black
		color.RGBA{R: 255, G: 0, B: 0},     // Red
		color.RGBA{R: 255, G: 165, B: 0},   // Orange
		color.RGBA{R: 255, G: 255, B: 0},   // Yellow
		color.RGBA{R: 0, G: 128, B: 0},     // Green
		color.RGBA{R: 0, G: 0, B: 255},     // Blue
		color.RGBA{R: 75, G: 0, B: 130},    // Indigo
		color.RGBA{R: 238, G: 130, B: 238}, // Violet
	)

	window := explorer.Window[float64]{
		Rect: geometry.Rect[float64]{
			Min: geometry.Point[float64]{X: xMin, Y: yMin},
			Max: geometry.Point[float64]{X: xMax, Y: yMax},
		},
	}

	image := explorer.NewGraph(
		screenWidth, screenHeight,
		iterator,
		&pallete,
	)

	// Create a new Game instance
	game := explorer.NewExplorer(
		screen,
		window,
		image,
	)

	// Set the window to full screen
	ebiten.SetFullscreen(true)

	// Run the game
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
