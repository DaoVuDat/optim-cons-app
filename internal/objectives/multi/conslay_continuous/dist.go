package conslay_continuous

import "math"

func Distance2D(a, b Coordinate) float64 {
	x := a.X - b.X
	y := a.Y - b.Y

	return math.Sqrt(x*x + y*y)
}
