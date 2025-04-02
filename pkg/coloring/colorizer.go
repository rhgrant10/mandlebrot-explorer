package coloring

import (
	"image/color"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/fractal"
)

type Colorizer interface {
	GetColor(result fractal.IterationResult) color.RGBA
	SetResolution(r int)
}
