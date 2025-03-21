package util

func MinWithIdx[T int | float64](arr []T) (T, int) {
	minIdx := 0
	var minVal T = arr[0]

	for i := 0; i < len(arr); i++ {
		if arr[i] < arr[minIdx] {
			minIdx = i
			minVal = arr[i]
		}
	}

	return minVal, minIdx
}

func MaxWithIdx[T int | float64](arr []T) (T, int) {
	maxIdx := 0
	var maxVal T = arr[0]

	for i := 0; i < len(arr); i++ {
		if arr[i] > arr[maxIdx] {
			maxIdx = i
			maxVal = arr[i]
		}
	}

	return maxVal, maxIdx
}
