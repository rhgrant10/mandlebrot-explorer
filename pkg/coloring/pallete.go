package coloring

import (
	"image/color"
	"math"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/fractal"
)

func NewPallete(resolution int, outColor color.RGBA, inColor color.RGBA, inColors ...color.RGBA) Pallete {
	return Pallete{
		inColors:   append([]color.RGBA{inColor}, inColors...),
		outColor:   outColor,
		resolution: resolution,
	}
}

type Pallete struct {
	inColors   []color.RGBA
	outColor   color.RGBA
	resolution int
}

func (p *Pallete) GetColor(result fractal.IterationResult) color.RGBA {
	clampedIterCount := math.Min(float64(result.IterCount), float64(p.resolution))
	ratio := clampedIterCount / float64(p.resolution+1)
	t := ratio * float64(len(p.inColors))
	si := int(t) // truncate
	ei := si + 1
	if ei >= int(len(p.inColors)) {
		return p.outColor
	}
	return p.Blend(p.inColors[si], p.inColors[ei], t-float64(si))
}

func (p *Pallete) SetResolution(r int) {
	p.resolution = r
}

func (p *Pallete) Blend(a, b color.RGBA, t float64) color.RGBA {
	return color.RGBA{
		R: float64toUnit8(float64(a.R) + t*(float64(b.R)-float64(a.R))),
		G: float64toUnit8(float64(a.G) + t*(float64(b.G)-float64(a.G))),
		B: float64toUnit8(float64(a.B) + t*(float64(b.B)-float64(a.B))),
		A: 255,
	}
}
