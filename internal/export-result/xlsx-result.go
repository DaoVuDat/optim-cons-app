package export_result

import "golang-moaha-construction/internal/algorithms"

type Summary struct {
	AlgorithmInfo   any
	ConstraintsInfo any
	ProblemInfo     any
	ObjectivesInfo  any
}

type ResultSummary struct {
	Idx int
}

type Options struct {
	Summary  Summary
	Results  algorithms.Result
	FilePath string
}

func WriteXlsxResult(option Options) error {

	return nil
}
