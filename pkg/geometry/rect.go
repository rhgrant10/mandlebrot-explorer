package geometry

func SizedRect[T Number](width, height T) Rect[T] {
	return Rect[T]{
		Min: Point[T]{},
		Max: Point[T]{X: width, Y: height},
	}
}

type Rect[T Number] struct {
	Min Point[T]
	Max Point[T]
}

func (r *Rect[T]) Width() T {
	return r.Max.X - r.Min.X
}

func (r *Rect[T]) Height() T {
	return r.Max.Y - r.Min.Y
}

func (r *Rect[T]) AspectRatio() float64 {
	return float64(r.Width()) / float64(r.Height())
}

func (r *Rect[T]) Area() float64 {
	return float64(r.Width() * r.Height())
}

func (r *Rect[T]) Size() Point[T] {
	return Point[T]{X: r.Width(), Y: r.Height()}
}

func (r *Rect[T]) CenterPoint() Point[T] {
	return Point[T]{
		X: r.Min.X + r.Width()/2,
		Y: r.Min.Y + r.Height()/2,
	}
}

func (r *Rect[T]) Translate(d Point[T]) *Rect[T] {
	r.Min.Translate(d)
	r.Max.Translate(d)
	return r
}

func (r *Rect[T]) Scale(f float64) *Rect[T] {
	r.Min.Scale(f)
	r.Max.Scale(f)
	return r
}
