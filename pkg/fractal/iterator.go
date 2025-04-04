package fractal

import "math/cmplx"

type IterationResult struct {
	FinalPoint PointPair
	IterCount  int
}

type Iterator struct {
	Equation  Equation
	AbsLimit  float64
	IterLimit int
}

func (it *Iterator) Iterate(p PointPair) IterationResult {
	n := 0
	for ; n < it.IterLimit; n++ {
		if cmplx.Abs(p.Z) > it.AbsLimit {
			break
		}
		p = it.Equation(p)
	}
	return IterationResult{
		FinalPoint: p,
		IterCount:  n,
	}
}
