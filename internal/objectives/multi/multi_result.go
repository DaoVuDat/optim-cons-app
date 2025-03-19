package multi

import "golang-moaha-construction/internal/objectives"

type MultiResult struct {
	objectives.Result
	CrowdingDistance float64
	Dominated        bool
	Rank             []int
	DominationSet    []MultiResult
	DominatedCount   []int
}
