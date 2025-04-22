package constraints

import (
	"golang-moaha-construction/internal/data"
	"math"
	"testing"
)

var (
	largeFacilityNames []string = []string{"TF1", "TF2"}
	smallLocationNames []string = []string{"L2", "L3"}
)

func createInputLocation(caseNumber int) map[string]data.Location {
	locations := make(map[string]data.Location)

	switch caseNumber {
	case 0:
		locations = map[string]data.Location{
			"TF1": data.Location{
				Symbol:      "TF1",
				IsLocatedAt: "L1",
			},
			"TF2": data.Location{
				Symbol:      "TF2",
				IsLocatedAt: "L4",
			},

			"TF3": data.Location{
				Symbol:      "TF3",
				IsLocatedAt: "L2",
			},
			"TF4": data.Location{
				Symbol:      "TF4",
				IsLocatedAt: "L3",
			},
		}
	case 1:
		locations = map[string]data.Location{
			"TF1": data.Location{
				Symbol:      "TF1",
				IsLocatedAt: "L2",
			},
			"TF2": data.Location{
				Symbol:      "TF2",
				IsLocatedAt: "L3",
			},

			"TF3": data.Location{
				Symbol:      "TF3",
				IsLocatedAt: "L1",
			},
			"TF4": data.Location{
				Symbol:      "TF4",
				IsLocatedAt: "L4",
			},
		}
	case 2:
		locations = map[string]data.Location{
			"TF1": data.Location{
				Symbol:      "TF1",
				IsLocatedAt: "L1",
			},
			"TF2": data.Location{
				Symbol:      "TF2",
				IsLocatedAt: "L2",
			},

			"TF3": data.Location{
				Symbol:      "TF3",
				IsLocatedAt: "L3",
			},
			"TF4": data.Location{
				Symbol:      "TF4",
				IsLocatedAt: "L4",
			},
		}
	default:

	}

	return locations
}

func TestSizeConstraint_Eval(t *testing.T) {
	testTable := []struct {
		locations map[string]data.Location
		expected  float64
		name      string
	}{
		{
			locations: createInputLocation(0),
			expected:  0,
			name:      "all valid located facilities",
		},
		{
			locations: createInputLocation(1),
			expected:  2,
			name:      "two invalid located facilities",
		},
		{
			locations: createInputLocation(2),
			expected:  1,
			name:      "one invalid located facility",
		},
	}

	sizeConstraint := CreateSizeConstraint(
		smallLocationNames,
		largeFacilityNames,
		1,
		1,
	)

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			result := sizeConstraint.Eval(tt.locations)

			if math.Abs(tt.expected-result) > 1e-9 {
				t.Errorf("expected result to be %f, got %f", tt.expected, result)
			}
		})
	}
}
