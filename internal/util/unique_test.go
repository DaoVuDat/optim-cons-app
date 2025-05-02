package util

import (
	"reflect"
	"testing"
)

func TestUniqueInts(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "no duplicates",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "with duplicates",
			input:    []int{1, 2, 2, 3, 4, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "all duplicates",
			input:    []int{1, 1, 1, 1},
			expected: []int{1},
		},
		{
			name:     "negative numbers",
			input:    []int{-1, -2, -1, 0, 1, 2, 1},
			expected: []int{-1, -2, 0, 1, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqueInts(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqueInts() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqueFloats(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{
			name:     "empty slice",
			input:    []float64{},
			expected: []float64{},
		},
		{
			name:     "no duplicates",
			input:    []float64{1.1, 2.2, 3.3, 4.4, 5.5},
			expected: []float64{1.1, 2.2, 3.3, 4.4, 5.5},
		},
		{
			name:     "with duplicates",
			input:    []float64{1.1, 2.2, 2.2, 3.3, 4.4, 4.4, 5.5},
			expected: []float64{1.1, 2.2, 3.3, 4.4, 5.5},
		},
		{
			name:     "all duplicates",
			input:    []float64{1.1, 1.1, 1.1, 1.1},
			expected: []float64{1.1},
		},
		{
			name:     "negative numbers",
			input:    []float64{-1.1, -2.2, -1.1, 0.0, 1.1, 2.2, 1.1},
			expected: []float64{-1.1, -2.2, 0.0, 1.1, 2.2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqueFloats(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqueFloats() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestUniqueStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c", "d", "e"},
			expected: []string{"a", "b", "c", "d", "e"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "b", "c", "d", "d", "e"},
			expected: []string{"a", "b", "c", "d", "e"},
		},
		{
			name:     "all duplicates",
			input:    []string{"a", "a", "a", "a"},
			expected: []string{"a"},
		},
		{
			name:     "empty strings",
			input:    []string{"", "a", "", "b", "c"},
			expected: []string{"", "a", "b", "c"},
		},
		{
			name:     "case sensitivity",
			input:    []string{"a", "A", "b", "B", "a"},
			expected: []string{"a", "A", "b", "B"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UniqueStrings(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("UniqueStrings() = %v, want %v", result, tt.expected)
			}
		})
	}
}