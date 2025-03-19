package util

import "math/rand"

func RandN(dim int) []int {
	res := make([]int, dim)
	for i := 0; i < dim; i++ {
		res[i] = i
	}

	rand.Shuffle(dim, func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	return res
}
