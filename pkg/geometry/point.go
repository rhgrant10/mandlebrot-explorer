package geometry

type Number interface {
	int | float64
}

type Point[T Number] struct {
	X T
	Y T
}

func (p *Point[T]) AsComplex() complex128 {
	return complex(float64(p.X), float64(p.Y))
}

func (p *Point[T]) Translate(d Point[T]) *Point[T] {
	p.X += d.X
	p.Y += d.Y
	return p
}

func (r *Point[T]) Scale(f float64) *Point[T] {
	r.X = T(float64(r.X) * f)
	r.Y = T(float64(r.Y) * f)
	return r
}

func (r *Point[T]) Negate() *Point[T] {
	r.X = -r.X
	r.Y = -r.Y
	return r
}
