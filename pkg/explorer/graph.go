package explorer

import (
	"fmt"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/coloring"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/fractal"
	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"
)

func NewGraph(
	width, height int,
	iterator fractal.Iterator,
	colorizer coloring.Colorizer,
	concurrency int,
) Graph {
	return Graph{
		Iterator:    iterator,
		Colorizer:   colorizer,
		canvas:      NewCanvas(width, height),
		concurrency: concurrency,
	}
}

type Graph struct {
	Iterator    fractal.Iterator
	Colorizer   coloring.Colorizer
	isUpdating  bool
	canvas      Canvas
	concurrency int
}

func (im *Graph) Resolution() geometry.Point[int] {
	return im.canvas.Size()
}

func (im *Graph) Render(window Window[float64]) error {
	fmt.Printf("Rendering\n")
	im.isUpdating = true
	resolution := im.Resolution()
	wg := sync.WaitGroup{}
	start := time.Now()
	iterCount := 0
	sem := make(chan struct{}, im.concurrency)
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
				iterCount += result.IterCount
				color := im.Colorizer.GetColor(result)
				im.canvas.Set(x, y, color)
			}()
		}
	}
	wg.Wait()
	im.isUpdating = false
	elapsed := time.Since(start)
	fmt.Printf("Performed %e iterations per second\n", float64(iterCount)/elapsed.Seconds())
	return nil
}

func (im *Graph) IsRendering() bool {
	return im.isUpdating
}

func (im *Graph) DrawOn(image *ebiten.Image) {
	im.canvas.DrawImage(image)
}
