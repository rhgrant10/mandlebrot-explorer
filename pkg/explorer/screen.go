package explorer

import (
	"image/color"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"
)

const (
	bytesPerPixel = 4
)

func NewScreen(w, h int) Screen {
	return Screen{
		rect:   geometry.SizedRect(w, h),
		pixels: make([]byte, w*h*bytesPerPixel),
	}
}

type Screen struct {
	rect   geometry.Rect[int]
	pixels []byte
}

func (s *Screen) Width() int {
	return s.rect.Width()
}

func (s *Screen) Height() int {
	return s.rect.Height()
}

func (s *Screen) AspectRatio() float64 {
	return s.rect.AspectRatio()
}

func (s *Screen) Pixels() []byte {
	return s.pixels
}

func (s *Screen) Color(x, y int, c color.RGBA) {
	idx := (y*s.rect.Width() + x) * bytesPerPixel
	s.pixels[idx] = c.R
	s.pixels[idx+1] = c.G
	s.pixels[idx+2] = c.B
	s.pixels[idx+3] = c.A
}
