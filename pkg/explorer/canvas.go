package explorer

import (
	"image/color"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewCanvas(w, h int) Canvas {
	return Canvas{
		Rect:   geometry.SizedRect(w, h),
		pixels: make([]byte, w*h*bytesPerPixel),
		image:  ebiten.NewImage(w, h),
	}
}

type Canvas struct {
	geometry.Rect[int]
	Options *ebiten.DrawImageOptions
	pixels  []byte
	image   *ebiten.Image
}

func (c *Canvas) DrawImage(image *ebiten.Image) {
	c.image.WritePixels(c.pixels)
	image.DrawImage(c.image, c.Options)
}

func (c *Canvas) Set(x, y int, color color.RGBA) {
	idx := (y*c.Rect.Width() + x) * bytesPerPixel
	c.pixels[idx] = color.R
	c.pixels[idx+1] = color.G
	c.pixels[idx+2] = color.B
	c.pixels[idx+3] = color.A
}
