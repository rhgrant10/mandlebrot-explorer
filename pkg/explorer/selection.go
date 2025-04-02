package explorer

import (
	"image/color"

	"github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewSelection(size uint8) Selection {
	return Selection{displaySize: int(size), aspectRatio: 1.0}
}

type Selection struct {
	displaySize int
	aspectRatio float64
	start       *geometry.Point[int]
}

func (s *Selection) GetMouse() *geometry.Point[int] {
	mx, my := ebiten.CursorPosition()
	return &geometry.Point[int]{X: mx, Y: my}
}

func (s *Selection) Start() {
	s.start = s.GetMouse()
}

func (s *Selection) End() geometry.Rect[int] {
	rect := s.Get()
	s.Clear()
	return rect
}

func (s *Selection) Get() geometry.Rect[int] {
	mouse := *s.GetMouse()
	if !s.IsBeingMade() {
		return geometry.Rect[int]{Min: mouse, Max: mouse}
	}
	rect := geometry.Rect[int]{Min: *s.start, Max: mouse}
	if rect.Width() < s.displaySize {
		if rect.Max.X < rect.Min.X {
			rect.Max.X = rect.Min.X - s.displaySize
		} else {
			rect.Max.X = rect.Min.X + s.displaySize
		}
	}
	h := int(float64(rect.Width()) / s.aspectRatio)
	if rect.Max.X < rect.Min.X {
		h = -h
	}
	if rect.Max.Y < rect.Min.Y {
		h = -h
	}
	rect.Max.Y = rect.Min.Y + h
	return geometry.Rect[int]{
		Min: geometry.Point[int]{
			X: min(rect.Min.X, rect.Max.X),
			Y: min(rect.Min.Y, rect.Max.Y),
		},
		Max: geometry.Point[int]{
			X: max(rect.Min.X, rect.Max.X),
			Y: max(rect.Min.Y, rect.Max.Y),
		},
	}
}

func (s *Selection) Clear() {
	s.start = nil
}

func (s *Selection) IsBeingMade() bool {
	return s.start != nil
}

func (s *Selection) Render(screen *ebiten.Image, window Window[float64]) {
	size := screen.Bounds()
	s.aspectRatio = float64(size.Dx()) / float64(size.Dy())
	if s.start == nil {
		return
	}
	rect := s.Get()
	for _, y := range []int{rect.Min.Y, rect.Max.Y} {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			var v = 255 * uint8((x/s.displaySize)%2)
			c := color.RGBA{v, v, v, 255}
			screen.Set(x, y, c)
		}
	}
	for _, x := range []int{rect.Min.X, rect.Max.X} {
		for y := rect.Min.Y; y < rect.Max.Y; y++ {
			var v = 255 * uint8((y/s.displaySize)%2)
			c := color.RGBA{v, v, v, 255}
			screen.Set(x, y, c)
		}
	}
}
