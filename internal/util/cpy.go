package util

func CopyArray[T int | float64](src []T) []T {
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

func CopyMap[T int | float64](src map[string]T) map[string]T {
	dst := make(map[string]T, len(src))

	for k, v := range src {
		dst[k] = v
	}

	return dst
}
