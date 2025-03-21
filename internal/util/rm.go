package util

func Remove[T any](arr []T, idx int) []T {
	if idx == len(arr)-1 {
		return arr[:idx]
	} else if idx == 0 {
		return arr[1:]
	} else {
		return append(arr[:idx], arr[idx+1:]...)
	}
}
