package fractal

type Equation func(PointPair) PointPair

func Mandlebrot(p PointPair) PointPair {
	return PointPair{
		Z: p.Z*p.Z + p.C,
		C: p.C,
	}
}
