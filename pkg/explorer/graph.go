package explorer

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/coloring"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/fractal"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"
)

func NewGraph(
	width, height int,
	iterator fractal.Iterator,
	colorizer coloring.Colorizer,
) Graph {
	return Graph{
		Iterator:  iterator,
		Colorizer: colorizer,
		canvas:    NewCanvas(width, height),
	}
}

type Graph struct {
	Iterator   fractal.Iterator
	Colorizer  coloring.Colorizer
	isUpdating bool
	canvas     Canvas
}

func (im *Graph) Resolution() geometry.Point[int] {
	return im.canvas.Size()
}

func (im *Graph) Render(window Window[float64]) error {
	im.isUpdating = true
	resolution := im.Resolution()
	wg := sync.WaitGroup{}
	sem := make(chan struct{}, 120)
	for y := 0; y < resolution.Y; y++ {
		for x := 0; x < resolution.X; x++ {
			wg.Add(1)
			sem <- struct{}{}
			go func() {
				defer wg.Done()
				defer func() { <-sem }()
				c := window.Transform(x, y, resolution.X, resolution.Y)
				p := fractal.PointPair{
					Z: 0,
					C: c.AsComplex(),
				}
				result := im.Iterator.Iterate(p)
				color := im.Colorizer.GetColor(result)
				im.canvas.Set(x, y, color)
			}()
		}
	}
	wg.Wait()
	im.isUpdating = false
	return nil
}

func (im *Graph) IsRendering() bool {
	return im.isUpdating
}

func (im *Graph) DrawOn(image *ebiten.Image) {
	im.canvas.DrawImage(image)
}
