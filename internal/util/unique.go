package util

// UniqueInts returns a new slice containing only the unique integer values from the input slice.
// It preserves the order of the first occurrence of each value.
func UniqueInts(input []int) []int {
	if len(input) == 0 {
		return []int{}
	}

	seen := make(map[int]struct{})
	result := make([]int, 0, len(input))

	for _, val := range input {
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}

// UniqueFloats returns a new slice containing only the unique float64 values from the input slice.
// It preserves the order of the first occurrence of each value.
func UniqueFloats(input []float64) []float64 {
	if len(input) == 0 {
		return []float64{}
	}

	seen := make(map[float64]struct{})
	result := make([]float64, 0, len(input))

	for _, val := range input {
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}

// UniqueStrings returns a new slice containing only the unique string values from the input slice.
// It preserves the order of the first occurrence of each value.
func UniqueStrings(input []string) []string {
	if len(input) == 0 {
		return []string{}
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(input))

	for _, val := range input {
		if _, exists := seen[val]; !exists {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}

	return result
}
