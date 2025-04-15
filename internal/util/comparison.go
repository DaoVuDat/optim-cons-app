package util

import (
	"regexp"
	"strconv"
)

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

var re = regexp.MustCompile(`\D*(\d+)`)

func ExtractNumber(s string) int {
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0
	}
	num, _ := strconv.Atoi(matches[1])
	return num
}
