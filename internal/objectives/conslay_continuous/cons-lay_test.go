package conslay_continuous

import "testing"

func TestReadLocationsFromFile(t *testing.T) {
	locs, fixedLocs, nonFixedLocs, err := ReadLocationsFromFile("../../../data/conslay/mini/locations.xlsx")

	if err != nil {
		t.Errorf("expected to read phases from file, got %s", err)
	}

	if len(locs) != 4 {
		t.Errorf("expected to read 4 locations, got %d", len(locs))
	}

	if len(fixedLocs) != 1 {
		t.Errorf("expected to read 1 fixed locations, got %d", len(fixedLocs))
	}

	if len(nonFixedLocs) != 3 {
		t.Errorf("expected to read 3 non-fixed locations, got %d", len(nonFixedLocs))
	}
}
