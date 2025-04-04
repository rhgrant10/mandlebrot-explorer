package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/coloring"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/explorer"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/fractal"
)

const (
	xMin        float64 = -3.0
	xMax        float64 = 1.2
	yMin        float64 = -1.2
	yMax        float64 = 1.2
	absLimit    float64 = 2.0
	iterLimit   int     = 1000
	concurrency int     = 1000
)

func main() {
	// Set the window to full screen
	ebiten.SetFullscreen(true)
	monitor := ebiten.Monitor()
	screenWidth, screenHeight := monitor.Size()

	iterator := fractal.Iterator{
		Equation:  fractal.Mandlebrot,
		AbsLimit:  absLimit,
		IterLimit: iterLimit,
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

	window := explorer.NewWindow(xMin, yMin, xMax, yMax)
	window.ZoomTo(-0.7463, 0.1102, 100.0)

	graph := explorer.NewGraph(
		screenWidth, screenHeight,
		iterator,
		&pallete,
		concurrency,
	)

	// Create the explorerGame and run it
	explorerGame := explorer.NewGame(
		window,
		graph,
	)

	// Run the game
	if err := ebiten.RunGame(&explorerGame); err != nil {
		log.Fatal(err)
	}
}
