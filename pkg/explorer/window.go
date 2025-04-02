package explorer

import "github.com/rhgrant10/mandlebrot-explorer/pkg/geometry"

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
