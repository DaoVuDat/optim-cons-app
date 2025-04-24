package util

import (
	"cmp"
	"slices"
)

func SortWithIdx[T cmp.Ordered](values []T) ([]T, []int) {
	type idxValue struct {
		Value T
		Idx   int
	}

	toBeSorted := make([]idxValue, len(values))

	for i, v := range values {
		toBeSorted[i].Value = v
		toBeSorted[i].Idx = i
	}

	slices.SortStableFunc(toBeSorted, func(i, j idxValue) int {
		return cmp.Compare(i.Value, j.Value)
	})

	resT := make([]T, len(values))
	resIdx := make([]int, len(values))

	for i, v := range toBeSorted {
		resT[i] = v.Value
		resIdx[i] = v.Idx
	}

	return resT, resIdx
}
