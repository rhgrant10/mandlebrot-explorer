package explorer

import "github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"

func NewWindow[T geometry.Number](xMin, yMin, xMax, yMax T) Window[T] {
	return Window[T]{
		Rect: geometry.Rect[T]{
			Min: geometry.Point[T]{X: xMin, Y: yMin},
			Max: geometry.Point[T]{X: xMax, Y: yMax},
		},
	}
}

func NewWindowInto[T geometry.Number](x, y T, zoom float64, original Window[T]) Window[T] {
	w := T(float64(original.Width()) / zoom)
	h := T(float64(original.Height()) / zoom)
	return NewWindow(x-w/2, y-h/2, x+w/2, y+h/2)
}

type Window[T geometry.Number] struct {
	geometry.Rect[T]
}

func (win *Window[T]) Transform(sx, sy, sw, sh int) geometry.Point[T] {
	xr := float64(sx) / float64(sw)
	yr := float64(sy) / float64(sh)
	return geometry.Point[T]{
		X: win.Min.X + T(xr*float64(win.Width())),
		Y: win.Min.Y + T(yr*float64(win.Height())),
	}
}

func (win *Window[T]) ZoomTo(x, y T, zoom float64) {
	w := T(float64(win.Width()) / zoom)
	h := T(float64(win.Height()) / zoom)
	win.Rect.Min.X = x - w/2
	win.Rect.Min.Y = y - h/2
	win.Rect.Max.X = x + w/2
	win.Rect.Max.Y = y + h/2
}
