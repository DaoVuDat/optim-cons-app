package util

import "math"

func RoundTo(n float64, decimals uint32) float64 {
	return math.Round(n*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

func RoundToGrid(n float64, gridSize int) float64 {

	return math.Round(n/float64(gridSize)) * float64(gridSize)
}
