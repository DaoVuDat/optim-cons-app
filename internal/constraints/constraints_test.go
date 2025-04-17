package constraints

import (
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/util"
	"math"
	"testing"
)

func TestIsCoverRangeOfCrane(t *testing.T) {
	tests := []struct {
		name      string
		crane     data.Crane
		buildings []data.Location
		expected  bool
		invalid   float64
	}{
		{
			name: "crane fully covers one building",
			crane: data.Crane{
				Location: data.Location{
					Coordinate: data.Coordinate{X: 0, Y: 0},
				},
				Radius: 10,
			},
			buildings: []data.Location{
				{
					Coordinate: data.Coordinate{X: 2, Y: 2},
					Length:     4,
					Width:      4,
				},
			},
			expected: true,
			invalid:  0,
		},
		{
			name: "crane partially covers a building",
			crane: data.Crane{
				Location: data.Location{
					Coordinate: data.Coordinate{X: 0, Y: 0},
				},
				Radius: 8,
			},
			buildings: []data.Location{
				{
					Coordinate: data.Coordinate{X: 6, Y: 6},
					Length:     2,
					Width:      2,
				},
			},
			expected: false,
			invalid: (data.Distance2D(data.Coordinate{X: 6 + 1, Y: 6 + 1}, data.Coordinate{X: 0, Y: 0}) - 8) +
				(data.Distance2D(data.Coordinate{X: 6 + 1, Y: 6 - 1}, data.Coordinate{X: 0, Y: 0}) - 8) +
				(data.Distance2D(data.Coordinate{X: 6 - 1, Y: 6 + 1}, data.Coordinate{X: 0, Y: 0}) - 8),
		},
		{
			name: "crane doesn't cover any buildings",
			crane: data.Crane{
				Location: data.Location{
					Coordinate: data.Coordinate{X: 0, Y: 0},
				},
				Radius: 3,
			},
			buildings: []data.Location{
				{
					Coordinate: data.Coordinate{X: 10, Y: 10},
					Length:     4,
					Width:      4,
				},
			},
			expected: false,
			invalid: (data.Distance2D(data.Coordinate{X: 10 + 2, Y: 10 + 2}, data.Coordinate{X: 0, Y: 0}) - 3) +
				(data.Distance2D(data.Coordinate{X: 10 - 2, Y: 10 - 2}, data.Coordinate{X: 0, Y: 0}) - 3) +
				(data.Distance2D(data.Coordinate{X: 10 - 2, Y: 10 + 2}, data.Coordinate{X: 0, Y: 0}) - 3) +
				(data.Distance2D(data.Coordinate{X: 10 + 2, Y: 10 - 2}, data.Coordinate{X: 0, Y: 0}) - 3),
		},
		{
			name: "crane fully covers multiple buildings",
			crane: data.Crane{
				Location: data.Location{
					Coordinate: data.Coordinate{X: 0, Y: 0},
				},
				Radius: 15,
			},
			buildings: []data.Location{
				{
					Coordinate: data.Coordinate{X: 5, Y: 5},
					Length:     4,
					Width:      4,
				},
				{
					Coordinate: data.Coordinate{X: -5, Y: -5},
					Length:     4,
					Width:      4,
				},
			},
			expected: true,
			invalid:  0,
		},
		{
			name: "crane partially covers multiple buildings",
			crane: data.Crane{
				Location: data.Location{
					Coordinate: data.Coordinate{X: 0, Y: 0},
				},
				Radius: 5,
			},
			buildings: []data.Location{
				{
					Coordinate: data.Coordinate{X: 1, Y: 1},
					Length:     2,
					Width:      2,
				},
				{
					Coordinate: data.Coordinate{X: -6, Y: -6},
					Length:     2,
					Width:      2,
				},
			},
			expected: false,
			invalid: (data.Distance2D(data.Coordinate{X: -6 + 1, Y: -6 + 1}, data.Coordinate{X: 0, Y: 0}) - 5) +
				(data.Distance2D(data.Coordinate{X: -6 + 1, Y: -6 - 1}, data.Coordinate{X: 0, Y: 0}) - 5) +
				(data.Distance2D(data.Coordinate{X: -6 - 1, Y: -6 + 1}, data.Coordinate{X: 0, Y: 0}) - 5) +
				(data.Distance2D(data.Coordinate{X: -6 - 1, Y: -6 - 1}, data.Coordinate{X: 0, Y: 0}) - 5),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, invalidAmount := IsCoverRangeOfCrane(tt.crane, tt.buildings)

			if result != tt.expected || math.Abs(invalidAmount-tt.invalid) > 1e-9 {
				t.Errorf("IsCoverRangeOfCrane(%v, %v) = (%v, %f); want (%v, %f)",
					tt.crane, tt.buildings, result, invalidAmount, tt.expected, tt.invalid)
			}
		})
	}
}

func TestOverlapConstraintsFalse(t *testing.T) {

	testTable := []struct {
		name                     string
		b1                       data.Location
		b2                       data.Location
		expectedOverlappedAmount float64
		expectedIsOverlapped     bool
	}{
		{
			name:                     "Not overlap at top right",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 16, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at top",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at top left",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 4, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at left",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 4, Y: 10}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at bottom left",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 4, Y: 4}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at bottom",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 4}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at right",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 16}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
		{
			name:                     "Not overlap at bottom right",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 16, Y: 4}, Length: 2, Width: 2},
			expectedIsOverlapped:     false,
			expectedOverlappedAmount: 0,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOverlapped, overlappedAmount := IsOverlapped(tt.b1, tt.b2)
			if tt.expectedIsOverlapped != isOverlapped {
				t.Errorf("expected %t, get %t", tt.expectedIsOverlapped, isOverlapped)
			}

			if overlappedAmount != tt.expectedOverlappedAmount {
				t.Errorf("expected %f, get %f", tt.expectedOverlappedAmount, overlappedAmount)
			}
		})
	}

}

func TestOverlapConstraintsTrue(t *testing.T) {

	testTable := []struct {
		name                     string
		b1                       data.Location
		b2                       data.Location
		expectedOverlappedAmount float64
		expectedIsOverlapped     bool
	}{
		{
			name:                     "Overlap at top right",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 16, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
		{
			name:                     "Overlap at top",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Overlap at top left",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 4, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
		{
			name:                     "Overlap at left",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 4, Y: 10}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Not overlap at bottom left",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 4, Y: 4}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
		{
			name:                     "Overlap at bottom",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 4}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Not overlap at right",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 16}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 7,
		},
		{
			name:                     "Overlap at bottom right",
			b1:                       data.Location{Coordinate: data.Coordinate{X: 10, Y: 10}, Length: 10, Width: 10},
			b2:                       data.Location{Coordinate: data.Coordinate{X: 16, Y: 4}, Length: 3, Width: 3},
			expectedIsOverlapped:     true,
			expectedOverlappedAmount: 1,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOverlapped, overlappedAmount := IsOverlapped(tt.b1, tt.b2)
			if tt.expectedIsOverlapped != isOverlapped {
				t.Errorf("expected %t, get %t", tt.expectedIsOverlapped, isOverlapped)
			}

			if overlappedAmount != tt.expectedOverlappedAmount {
				t.Errorf("expected %f, get %f", tt.expectedOverlappedAmount, overlappedAmount)
			}
		})
	}

}

func TestIsOutOfBoundTrue(t *testing.T) {
	testTable := []struct {
		name                       string
		b                          data.Location
		minL                       float64
		minW                       float64
		maxL                       float64
		maxW                       float64
		expectedIsOutOfBoundAmount float64
		expectedIsOutOfBound       bool
	}{
		{
			name:                       "Out Of Bound at left",
			b:                          data.Location{Coordinate: data.Coordinate{X: 4, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at top left",
			b:                          data.Location{Coordinate: data.Coordinate{X: 4, Y: 91}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
		{
			name:                       "Out Of Bound at top",
			b:                          data.Location{Coordinate: data.Coordinate{X: 10, Y: 91}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at top right",
			b:                          data.Location{Coordinate: data.Coordinate{X: 116, Y: 91}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
		{
			name:                       "Out Of Bound at right",
			b:                          data.Location{Coordinate: data.Coordinate{X: 116, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at bottom right",
			b:                          data.Location{Coordinate: data.Coordinate{X: 116, Y: 4}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
		{
			name:                       "Out Of Bound at bottom",
			b:                          data.Location{Coordinate: data.Coordinate{X: 10, Y: 4}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 1,
		},
		{
			name:                       "Out Of Bound at bottom left",
			b:                          data.Location{Coordinate: data.Coordinate{X: 4, Y: 4}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       true,
			expectedIsOutOfBoundAmount: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOutOfBound, outOfBoundAmount := IsOutOfBound(tt.minL, tt.maxL, tt.minW, tt.maxW, tt.b)
			if tt.expectedIsOutOfBound != isOutOfBound {
				t.Errorf("expected %t, get %t", tt.expectedIsOutOfBound, isOutOfBound)
			}

			if outOfBoundAmount != tt.expectedIsOutOfBoundAmount {
				t.Errorf("expected %f, get %f", tt.expectedIsOutOfBoundAmount, outOfBoundAmount)
			}
		})
	}
}

func TestIsOutOfBoundFalse(t *testing.T) {
	testTable := []struct {
		name                       string
		b                          data.Location
		minL                       float64
		minW                       float64
		maxL                       float64
		maxW                       float64
		expectedIsOutOfBoundAmount float64
		expectedIsOutOfBound       bool
	}{
		{
			name:                       "Not out Of Bound at left",
			b:                          data.Location{Coordinate: data.Coordinate{X: 5, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at top left",
			b:                          data.Location{Coordinate: data.Coordinate{X: 5, Y: 90}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at top",
			b:                          data.Location{Coordinate: data.Coordinate{X: 10, Y: 90}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at top right",
			b:                          data.Location{Coordinate: data.Coordinate{X: 115, Y: 90}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at right",
			b:                          data.Location{Coordinate: data.Coordinate{X: 115, Y: 10}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at bottom right",
			b:                          data.Location{Coordinate: data.Coordinate{X: 115, Y: 5}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at bottom",
			b:                          data.Location{Coordinate: data.Coordinate{X: 10, Y: 5}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
		{
			name:                       "Not out Of Bound at bottom left",
			b:                          data.Location{Coordinate: data.Coordinate{X: 5, Y: 5}, Length: 10, Width: 10},
			minL:                       0,
			minW:                       0,
			maxL:                       120,
			maxW:                       95,
			expectedIsOutOfBound:       false,
			expectedIsOutOfBoundAmount: 0,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			isOutOfBound, outOfBoundAmount := IsOutOfBound(tt.minL, tt.maxL, tt.minW, tt.maxW, tt.b)
			if tt.expectedIsOutOfBound != isOutOfBound {
				t.Errorf("expected %t, get %t", tt.expectedIsOutOfBound, isOutOfBound)
			}

			if outOfBoundAmount != tt.expectedIsOutOfBoundAmount {
				t.Errorf("expected %f, get %f", tt.expectedIsOutOfBoundAmount, outOfBoundAmount)
			}
		})
	}
}

func CreateInputLocation(feasible bool) map[string]data.Location {
	//input := []float64{
	//	41.8984994960664, 16.0789402087708, 0.110667114989373,
	//	40.3018087791441, 6.54792941746860, 0.632349545372518,
	//	43.0538143534362, 81.9094405733262, 0.0819942691480144,
	//	49.7610024122533, 57.2431089610756, 0.845980667765488,
	//	81.3261517318950, 41.8071744013404, 0.274381079260441,
	//	101.490485902993, 11.1774705429599, 0.566086701250327,
	//	43.6721024451348, 9.42961655802787, 0.875812362483772,
	//	48.5660821450197, 54.0455646998945, 0.207017324531373,
	//	76.4180335261621, 41.8787903034157, 0.298433966901801,
	//	108.028665273018, 36.9845037867221, 0.712680757776487,
	//	63.8375790712078, 91.2094786904858, 0.746326249229960,
	//	110.579901257085, 14.7798350113455, 0.122726427684860,
	//} // feasible

	//input := []float64{
	//	60.2117021287296, 70.6492302951395, 0.840088279889616,
	//	70.8063495854695, 20.8484429816236, 0.930459707285301,
	//	16.0620484334653, 11.2728104676465, 0.256779152503716,
	//	66.8271432818314, 79.6322545348163, 0.994330695133932,
	//	97.4966707756999, 13.6022897445312, 0.142998150825580,
	//	1.51514052087480, 80.0719162226562, 0.0752335982572160,
	//	36.7111938858722, 11.9943304146193, 0.608734625270310,
	//	48.6579882916605, 61.3097045687976, 0.0711765983287040,
	//	7.04418109569817, 87.8287852386491, 0.736732593589075,
	//	42.2938936308434, 65.2135456501644, 0.156252552003082,
	//	21.9910306881593, 85.0044587159497, 0.464587709048142,
	//	23.2731387385839, 90.0053928795082, 0.790347360934255,
	//} // infeasible

	locations := make(map[string]data.Location)
	if feasible {
		locations = map[string]data.Location{
			"TF1": data.Location{
				Coordinate: data.Coordinate{X: 41.8984994960664, Y: 16.0789402087708},
				Rotation:   false,
				Length:     12,
				Width:      5,
				IsFixed:    false,
				Symbol:     "TF1",
				Name:       "Pile storage yard #1",
			},
			"TF2": data.Location{
				Coordinate: data.Coordinate{X: 40.3018087791441, Y: 6.54792941746860},
				Rotation:   true,
				Length:     5,
				Width:      12,
				IsFixed:    false,
				Symbol:     "TF2",
				Name:       "Pile storage yard #2",
			},
			"TF3": data.Location{
				Coordinate: data.Coordinate{X: 43.0538143534362, Y: 81.9094405733262},
				Rotation:   false,
				Length:     8,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF3",
				Name:       "Site office",
			},
			"TF4": data.Location{
				Coordinate: data.Coordinate{X: 49.7610024122533, Y: 57.2431089610756},
				Rotation:   true,
				Length:     7,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF4",
				Name:       "Rebar process yard",
			},
			"TF5": data.Location{
				Coordinate: data.Coordinate{X: 81.3261517318950, Y: 41.8071744013404},
				Rotation:   false,
				Length:     14,
				Width:      7,
				IsFixed:    false,
				Symbol:     "TF5",
				Name:       "Formwork process yard",
			},
			"TF6": data.Location{
				Coordinate: data.Coordinate{X: 101.490485902993, Y: 11.1774705429599},
				Rotation:   true,
				Length:     4,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF6",
				Name:       "Electrician hut",
			},
			"TF7": data.Location{
				Coordinate: data.Coordinate{X: 43.6721024451348, Y: 9.42961655802787},
				Rotation:   true,
				Length:     12,
				Width:      10,
				IsFixed:    false,
				Symbol:     "TF7",
				Name:       "Ready-mix concrete area",
			},
			"TF8": data.Location{
				Coordinate: data.Coordinate{X: 48.5660821450197, Y: 54.0455646998945},
				Rotation:   false,
				Length:     12,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF8",
				Name:       "Prefabricated components yard #1 (slab, staircase)",
			},
			"TF9": data.Location{
				Coordinate: data.Coordinate{X: 76.4180335261621, Y: 41.8787903034157},
				Rotation:   false,
				Length:     12,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF9",
				Name:       "Prefabricated components yard #2 (beam, column)",
			},
			"TF10": data.Location{
				Coordinate: data.Coordinate{X: 108.028665273018, Y: 36.9845037867221},
				Rotation:   true,
				Length:     6,
				Width:      10,
				IsFixed:    false,
				Symbol:     "TF10",
				Name:       "Prefabricated components yard #3 (external wall)",
			},
			"TF11": data.Location{
				Coordinate: data.Coordinate{X: 63.8375790712078, Y: 91.2094786904858},
				Rotation:   true,
				Length:     4,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF11",
				Name:       "Dangerous goods warehouse",
			},
			"TF12": data.Location{
				Coordinate: data.Coordinate{X: 110.579901257085, Y: 14.7798350113455},
				Rotation:   false,
				Length:     8,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF12",
				Name:       "ME storage warehouse",
			},
			"TF13": data.Location{
				Coordinate: data.Coordinate{X: 72, Y: 22},
				Rotation:   false,
				Length:     52,
				Width:      24,
				IsFixed:    true,
				Symbol:     "TF13",
				Name:       "Building",
			},
			"TF14": {
				Coordinate: data.Coordinate{X: 72, Y: 36},
				Rotation:   false,
				Length:     2,
				Width:      2,
				IsFixed:    true,
				Symbol:     "TF14",
				Name:       "Tower crane",
			},
			"TF15": {
				Coordinate: data.Coordinate{X: 43, Y: 22},
				Rotation:   false,
				Length:     3,
				Width:      3,
				IsFixed:    true,
				Symbol:     "TF15",
				Name:       "Hoist",
			},
			"TF16": {
				Coordinate: data.Coordinate{X: 5, Y: 92},
				Rotation:   false,
				Length:     3,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF16",
				Name:       "Security room",
			},
		}
	} else {
		locations = map[string]data.Location{
			"TF1": data.Location{
				Coordinate: data.Coordinate{X: 60.2117021287296, Y: 70.6492302951395},
				Rotation:   true,
				Length:     5,
				Width:      12,
				IsFixed:    false,
				Symbol:     "TF1",
				Name:       "Pile storage yard #1",
			},
			"TF2": data.Location{
				Coordinate: data.Coordinate{X: 70.8063495854695, Y: 20.8484429816236},
				Rotation:   true,
				Length:     5,
				Width:      12,
				IsFixed:    false,
				Symbol:     "TF2",
				Name:       "Pile storage yard #2",
			},
			"TF3": data.Location{
				Coordinate: data.Coordinate{X: 16.0620484334653, Y: 11.2728104676465},
				Rotation:   false,
				Length:     8,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF3",
				Name:       "Site office",
			},
			"TF4": data.Location{
				Coordinate: data.Coordinate{X: 66.8271432818314, Y: 79.6322545348163},
				Rotation:   true,
				Length:     7,
				Width:      14,
				IsFixed:    false,
				Symbol:     "TF4",
				Name:       "Rebar process yard",
			},
			"TF5": data.Location{
				Coordinate: data.Coordinate{X: 97.4966707756999, Y: 13.6022897445312},
				Rotation:   false,
				Length:     14,
				Width:      7,
				IsFixed:    false,
				Symbol:     "TF5",
				Name:       "Formwork process yard",
			},
			"TF6": data.Location{
				Coordinate: data.Coordinate{X: 1.51514052087480, Y: 80.0719162226562},
				Rotation:   false,
				Length:     4,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF6",
				Name:       "Electrician hut",
			},
			"TF7": data.Location{
				Coordinate: data.Coordinate{X: 36.7111938858722, Y: 11.9943304146193},
				Rotation:   true,
				Length:     12,
				Width:      10,
				IsFixed:    false,
				Symbol:     "TF7",
				Name:       "Ready-mix concrete area",
			},
			"TF8": data.Location{
				Coordinate: data.Coordinate{X: 48.6579882916605, Y: 61.3097045687976},
				Rotation:   false,
				Length:     12,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF8",
				Name:       "Prefabricated components yard #1 (slab, staircase)",
			},
			"TF9": data.Location{
				Coordinate: data.Coordinate{X: 7.04418109569817, Y: 87.8287852386491},
				Rotation:   true,
				Length:     6,
				Width:      12,
				IsFixed:    false,
				Symbol:     "TF9",
				Name:       "Prefabricated components yard #2 (beam, column)",
			},
			"TF10": data.Location{
				Coordinate: data.Coordinate{X: 42.2938936308434, Y: 65.2135456501644},
				Rotation:   false,
				Length:     10,
				Width:      6,
				IsFixed:    false,
				Symbol:     "TF10",
				Name:       "Prefabricated components yard #3 (external wall)",
			},
			"TF11": data.Location{
				Coordinate: data.Coordinate{X: 21.9910306881593, Y: 85.0044587159497},
				Rotation:   false,
				Length:     6,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF11",
				Name:       "Dangerous goods warehouse",
			},
			"TF12": data.Location{
				Coordinate: data.Coordinate{X: 23.2731387385839, Y: 90.0053928795082},
				Rotation:   true,
				Length:     6,
				Width:      8,
				IsFixed:    false,
				Symbol:     "TF12",
				Name:       "ME storage warehouse",
			},
			"TF13": data.Location{
				Coordinate: data.Coordinate{X: 72, Y: 22},
				Rotation:   false,
				Length:     52,
				Width:      24,
				IsFixed:    true,
				Symbol:     "TF13",
				Name:       "Building",
			},
			"TF14": {
				Coordinate: data.Coordinate{X: 72, Y: 36},
				Rotation:   false,
				Length:     2,
				Width:      2,
				IsFixed:    true,
				Symbol:     "TF14",
				Name:       "Tower crane",
			},
			"TF15": {
				Coordinate: data.Coordinate{X: 43, Y: 22},
				Rotation:   false,
				Length:     3,
				Width:      3,
				IsFixed:    true,
				Symbol:     "TF15",
				Name:       "Hoist",
			},
			"TF16": {
				Coordinate: data.Coordinate{X: 5, Y: 92},
				Rotation:   false,
				Length:     3,
				Width:      4,
				IsFixed:    false,
				Symbol:     "TF16",
				Name:       "Security room",
			},
		}
	}

	return locations
}

func CreateInputPhases() [][]string {
	return [][]string{
		{"TF1", "TF2", "TF3", "TF6", "TF13", "TF16"},
		{"TF3", "TF4", "TF5", "TF6", "TF7", "TF13", "TF16"},
		{"TF3", "TF4", "TF5", "TF6", "TF7", "TF11", "TF13", "TF14", "TF15", "TF16"},
		{"TF3", "TF6", "TF8", "TF9", "TF11", "TF13", "TF14", "TF15", "TF16"},
		{"TF3", "TF6", "TF10", "TF11", "TF12", "TF13", "TF14", "TF15", "TF16"},
	}
}

func TestOutOfBoundsConstraint_Eval_Feasible(t *testing.T) {
	// Setup
	outOfBoundsConstraint := CreateOutOfBoundsConstraint(
		0,
		95,
		0,
		120,
		CreateInputPhases(),
		0,
		1,
	)

	locations := CreateInputLocation(true)

	penalty := outOfBoundsConstraint.Eval(locations)
	if math.Round(penalty) != 0 {
		t.Errorf("expected penalty to be 0, got %f", penalty)
	}

}

func TestOverlapConstraint_Eval_Feasible(t *testing.T) {
	overlapConstraint := CreateOverlapConstraint(
		CreateInputPhases(),
		0,
		1,
	)

	locations := CreateInputLocation(true)

	penalty := overlapConstraint.Eval(locations)
	if util.RoundTo(penalty, 2) != 26.68 {
		t.Errorf("expected penalty to be 26.68, got %f", penalty)
	}
}

func TestCoverRangeCraneConstraint_Eval_Feasible(t *testing.T) {
	locations := CreateInputLocation(true)
	coverRangeConstraint := CreateCoverRangeCraneConstraint(
		[]data.Crane{
			{
				Location:     locations["TF14"],
				BuildingName: []string{"TF4", "TF5", "TF8", "TF9", "TF10"},
				Radius:       40,
			},
		},
		CreateInputPhases(),
		0,
		1,
	)

	penalty := coverRangeConstraint.Eval(locations)
	if math.Round(penalty) != 0 {
		t.Errorf("expected penalty to be 0, got %f", penalty)
	}
}

func TestInclusiveZoneConstraint_Eval_Feasible(t *testing.T) {
	locations := CreateInputLocation(true)
	zoneConstraint := CreateInclusiveZoneConstraint(
		[]Zone{
			{
				Location:      locations["TF13"],
				BuildingNames: []string{"TF7"},
				Size:          20,
			},
			{
				Location:      locations["TF13"],
				BuildingNames: []string{"TF1", "TF2"},
				Size:          15,
			},
		},
		CreateInputPhases(),
		0,
		1,
	)

	penalty := zoneConstraint.Eval(locations)
	if math.Round(penalty) != 0 {
		t.Errorf("expected penalty to be 0, got %f", penalty)
	}
}

func TestOutOfBoundsConstraint_Eval_Infeasible(t *testing.T) {
	// Setup
	outOfBoundsConstraint := CreateOutOfBoundsConstraint(
		0,
		95,
		0,
		120,
		CreateInputPhases(),
		0,
		1,
	)

	locations := CreateInputLocation(false)

	penalty := outOfBoundsConstraint.Eval(locations)
	if util.RoundTo(penalty, 2) != 0.48 {
		t.Errorf("expected penalty to be 0.48, got %f", penalty)
	}

}

func TestOverlapConstraint_Eval_Infeasible(t *testing.T) {
	overlapConstraint := CreateOverlapConstraint(
		CreateInputPhases(),
		0,
		1,
	)

	locations := CreateInputLocation(false)

	penalty := overlapConstraint.Eval(locations)
	if util.RoundTo(penalty, 2) != 85.37 {
		t.Errorf("expected penalty to be 85.37, got %f", penalty)
	}
}

func TestCoverRangeCraneConstraint_Eval_Infeasible(t *testing.T) {
	locations := CreateInputLocation(false)
	coverRangeConstraint := CreateCoverRangeCraneConstraint(
		[]data.Crane{
			{
				Location:     locations["TF14"],
				BuildingName: []string{"TF4", "TF5", "TF8", "TF9", "TF10"},
				Radius:       40,
			},
		},
		CreateInputPhases(),
		0,
		1,
	)

	penalty := coverRangeConstraint.Eval(locations)
	if util.RoundTo(penalty, 2) != 208.81 {
		t.Errorf("expected penalty to be 208.81, got %f", penalty)
	}
}

func TestInclusiveZoneConstraint_Eval_Infeasible(t *testing.T) {
	locations := CreateInputLocation(false)
	zoneConstraint := CreateInclusiveZoneConstraint(
		[]Zone{
			{
				Location:      locations["TF13"],
				BuildingNames: []string{"TF7"},
				Size:          20,
			},
			{
				Location:      locations["TF13"],
				BuildingNames: []string{"TF1", "TF2"},
				Size:          15,
			},
		},
		CreateInputPhases(),
		0,
		1,
	)

	penalty := zoneConstraint.Eval(locations)
	if util.RoundTo(penalty, 2) != 27.65 {
		t.Errorf("expected penalty to be 27.65, got %f", penalty)
	}
}
