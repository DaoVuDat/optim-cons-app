package util

func CopyArray[T int | float64](src []T) []T {
	dst := make([]T, len(src))
	copy(dst, src)
	return dst
}
