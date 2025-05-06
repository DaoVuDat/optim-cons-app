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

// RouletteWheelSelection implements the roulette wheel selection algorithm.
// It takes a slice of probabilities as input and returns an index selected
// based on those probabilities. The higher the probability, the more likely
// the index is to be selected.
func RouletteWheelSelection(p []float64) int {
	if len(p) == 0 {
		return -1
	}

	// Calculate the sum of all probabilities
	sum := 0.0
	for _, prob := range p {
		sum += prob
	}

	// If sum is 0, return a random index
	if sum == 0 {
		return rand.Intn(len(p))
	}

	// Generate a random value between 0 and sum
	r := rand.Float64() * sum

	// Find the index corresponding to the random value
	currentSum := 0.0
	for i, prob := range p {
		currentSum += prob
		if r <= currentSum {
			return i
		}
	}

	// Fallback (should not reach here under normal circumstances)
	return len(p) - 1
}

func RandomSample(max int, n int) []int {
	// Ensure we don't try to sample more elements than available
	if n > max {
		n = max
	}

	// Create a slice of all indices
	indices := make([]int, max)
	for i := range indices {
		indices[i] = i
	}

	// Shuffle the indices
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	// Return the first n indices
	return indices[:n]
}
