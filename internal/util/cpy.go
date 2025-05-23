package util

import (
	"golang-moaha-construction/internal/data"
)

func CopyArray[T int | float64 | string | data.ObjectiveType](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}

func CopySliceOfSlice[T int | float64](src [][]T) [][]T {
	dst := make([][]T, len(src))
	for i := range src {
		dst[i] = make([]T, len(src[i]))
		copy(dst[i], src[i])
	}
	return dst
}

func CopyMap[T int | float64, K string | data.ConstraintType | data.ObjectiveType](src map[K]T) map[K]T {
	dst := make(map[K]T, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}
