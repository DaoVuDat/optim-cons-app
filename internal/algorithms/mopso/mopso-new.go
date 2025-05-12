package mopso

//
//import (
//	"golang-moaha-construction/internal/objectives"
//	"golang-moaha-construction/internal/util"
//	"math"
//	"math/rand"
//	"sort"
//)
//
//type MOPSOParams struct {
//	Np      int     // Number of particles
//	Nr      int     // Repository size
//	MaxGen  int     // Maximum generations
//	W       float64 // Inertia
//	C1      float64 // Personal confidence
//	C2      float64 // Swarm confidence
//	NGrid   int     // Number of hypercubes per dimension
//	MaxVel  float64 // Max velocity (percentage of search space)
//	UMut    float64 // Uniform mutation percentage
//}
//
//type MultiObj struct {
//	Fun      func([]float64) ([]float64, []float64) // returns (objectives, penalties)
//	NVar     int
//	VarMin   []float64
//	VarMax   []float64
//}
//
//type Repository struct {
//	Pos            [][]float64
//	PosFit         [][]float64
//	PosPen         [][]float64
//	HypercubeLims  [][]float64
//	GridIdx        []int
//	GridSubIdx     [][]int
//	Quality        [][2]float64 // [hypercube id, quality]
//}
//
//func MOPSO(params MOPSOParams, multiObj MultiObj) *Repository {
//	Np, Nr, MaxGen := params.Np, params.Nr, params.MaxGen
//	W, C1, C2 := params.W, params.C1, params.C2
//	ngrid, maxvel, u_mut := params.NGrid, params.MaxVel, params.UMut
//	fun, nVar := multiObj.Fun, multiObj.NVar
//	varMin, varMax := multiObj.VarMin, multiObj.VarMax
//
//	// Initialization
//	POS := make([][]float64, Np)
//	VEL := make([][]float64, Np)
//	for i := 0; i < Np; i++ {
//		POS[i] = make([]float64, nVar)
//		VEL[i] = make([]float64, nVar)
//		for j := 0; j < nVar; j++ {
//			POS[i][j] = varMin[j] + rand.Float64()*(varMax[j]-varMin[j])
//		}
//	}
//	maxVelVec := make([]float64, nVar)
//	for i := 0; i < nVar; i++ {
//		maxVelVec[i] = (varMax[i] - varMin[i]) * maxvel / 100.0
//	}
//
//	// Evaluate initial population
//	POSFit := make([][]float64, Np)
//	PENALTY := make([][]float64, Np)
//	for i := 0; i < Np; i++ {
//		fit, pen := fun(POS[i])
//		POSFit[i] = fit
//		PENALTY[i] = pen
//	}
//
//	PBEST := deepCopyMatrix(POS)
//	PBESTFit := deepCopyMatrix(POSFit)
//
//	// Initial repository
//	dominated := checkDomination(POSFit)
//	rep := &Repository{
//		Pos:    filterRows(POS, dominated, false),
//		PosFit: filterRows(POSFit, dominated, false),
//		PosPen: filterRows(PENALTY, dominated, false),
//	}
//	updateGrid(rep, ngrid)
//
//	gen := 1
//	for gen <= MaxGen {
//		// Select leader
//		h := selectLeader(rep)
//		leader := rep.Pos[h]
//
//		// Update speeds and positions
//		for i := 0; i < Np; i++ {
//			for j := 0; j < nVar; j++ {
//				r1 := rand.Float64()
//				r2 := rand.Float64()
//				VEL[i][j] = W*VEL[i][j] +
//					C1*r1*(PBEST[i][j]-POS[i][j]) +
//					C2*r2*(leader[j]-POS[i][j])
//				if math.Abs(VEL[i][j]) > maxVelVec[j] {
//					if VEL[i][j] > 0 {
//						VEL[i][j] = maxVelVec[j]
//					} else {
//						VEL[i][j] = -maxVelVec[j]
//					}
//				}
//				POS[i][j] += VEL[i][j]
//			}
//		}
//
//		// Mutation
//		mutation(POS, gen, MaxGen, Np, varMax, varMin, nVar, u_mut)
//
//		// Check boundaries
//		checkBoundaries(POS, VEL, maxVelVec, varMax, varMin)
//
//		// Evaluate
//		for i := 0; i < Np; i++ {
//			fit, pen := fun(POS[i])
//			POSFit[i] = fit
//			PENALTY[i] = pen
//		}
//
//		// Update repository
//		updateRepository(rep, POS, POSFit, PENALTY, ngrid)
//		if len(rep.Pos) > Nr {
//			deleteFromRepository(rep, len(rep.Pos)-Nr, ngrid)
//		}
//
//		// Update personal bests
//		for i := 0; i < Np; i++ {
//			if dominates(POSFit[i], PBESTFit[i]) {
//				PBESTFit[i] = util.CopyArray(POSFit[i])
//				PBEST[i] = util.CopyArray(POS[i])
//			} else if !dominates(PBESTFit[i], POSFit[i]) && rand.Float64() > 0.5 {
//				PBESTFit[i] = util.CopyArray(POSFit[i])
//				PBEST[i] = util.CopyArray(POS[i])
//			}
//		}
//		gen++
//	}
//	return rep
//}
//
//// --- Utility and helper functions ---
//
//func filterRows(mat [][]float64, mask []bool, keep bool) [][]float64 {
//	res := [][]float64{}
//	for i, row := range mat {
//		if mask[i] == keep {
//			res = append(res, util.CopyArray(row))
//		}
//	}
//	return res
//}
//
//func checkDomination(fitness [][]float64) []bool {
//	Np := len(fitness)
//	dom := make([]bool, Np)
//	for i := 0; i < Np; i++ {
//		for j := 0; j < Np; j++ {
//			if i != j && dominates(fitness[j], fitness[i]) {
//				dom[i] = true
//				break
//			}
//		}
//	}
//	return dom
//}
//
//func dominates(x, y []float64) bool {
//	allLE := true
//	anyLT := false
//	for i := range x {
//		if x[i] > y[i] {
//			allLE = false
//			break
//		}
//		if x[i] < y[i] {
//			anyLT = true
//		}
//	}
//	return allLE && anyLT
//}
//
//func updateGrid(rep *Repository, ngrid int) {
//	ndim := len(rep.PosFit[0])
//	rep.HypercubeLims = make([][]float64, ngrid+1)
//	for i := 0; i <= ngrid; i++ {
//		rep.HypercubeLims[i] = make([]float64, ndim)
//	}
//	for dim := 0; dim < ndim; dim++ {
//		minV, maxV := rep.PosFit[0][dim], rep.PosFit[0][dim]
//		for i := 1; i < len(rep.PosFit); i++ {
//			if rep.PosFit[i][dim] < minV {
//				minV = rep.PosFit[i][dim]
//			}
//			if rep.PosFit[i][dim] > maxV {
//				maxV = rep.PosFit[i][dim]
//			}
//		}
//		for i := 0; i <= ngrid; i++ {
//			rep.HypercubeLims[i][dim] = minV + float64(i)/float64(ngrid)*(maxV-minV)
//		}
//	}
//	// Assign grid indices
//	npar := len(rep.PosFit)
//	rep.GridIdx = make([]int, npar)
//	rep.GridSubIdx = make([][]int, npar)
//	for n := 0; n < npar; n++ {
//		rep.GridSubIdx[n] = make([]int, ndim)
//		idnames := make([]int, ndim)
//		for d := 0; d < ndim; d++ {
//			for k := 0; k < ngrid; k++ {
//				if rep.PosFit[n][d] <= rep.HypercubeLims[k+1][d] {
//					rep.GridSubIdx[n][d] = k + 1
//					break
//				}
//			}
//			if rep.GridSubIdx[n][d] == 0 {
//				rep.GridSubIdx[n][d] = 1
//			}
//			idnames[d] = rep.GridSubIdx[n][d] - 1
//		}
//		gridSize := make([]int, ndim)
//		for i := range gridSize {
//			gridSize[i] = ngrid
//		}
//		rep.GridIdx[n] = util.Sub2Index(gridSize, idnames...)
//	}
//	// Quality
//	count := map[int]int{}
//	for _, idx := range rep.GridIdx {
//		count[idx]++
//	}
//	rep.Quality = make([][2]float64, 0, len(count))
//	for idx, c := range count {
//		rep.Quality = append(rep.Quality, [2]float64{float64(idx), 10.0 / float64(c)})
//	}
//}
//
//func selectLeader(rep *Repository) int {
//	prob := make([]float64, len(rep.Quality))
//	prob[0] = rep.Quality[0][1]
//	for i := 1; i < len(rep.Quality); i++ {
//		prob[i] = prob[i-1] + rep.Quality[i][1]
//	}
//	r := rand.Float64() * prob[len(prob)-1]
//	selHyp := -1
//	for i, p := range prob {
//		if r <= p {
//			selHyp = int(rep.Quality[i][0])
//			break
//		}
//	}
//	indices := []int{}
//	for i, idx := range rep.GridIdx {
//		if idx == selHyp {
//			indices = append(indices, i)
//		}
//	}
//	if len(indices) == 0 {
//		return rand.Intn(len(rep.Pos))
//	}
//	return indices[rand.Intn(len(indices))]
//}
//
//func updateRepository(rep *Repository, POS, POSFit, POSPen [][]float64, ngrid int) {
//	dominated := checkDomination(POSFit)
//	rep.Pos = append(rep.Pos, filterRows(POS, dominated, false)...)
//	rep.PosFit = append(rep.PosFit, filterRows(POSFit, dominated, false)...)
//	rep.PosPen = append(rep.PosPen, filterRows(POSPen, dominated, false)...)
//	dominated2 := checkDomination(rep.PosFit)
//	rep.Pos = filterRows(rep.Pos, dominated2, false)
//	rep.PosFit = filterRows(rep.PosFit, dominated2, false)
//	rep.PosPen = filterRows(rep.PosPen, dominated2, false)
//	updateGrid(rep, ngrid)
//}
//
//func deleteFromRepository(rep *Repository, nExtra, ngrid int) {
//	crowding := make([]float64, len(rep.Pos))
//	for m := 0; m < len(rep.PosFit[0]); m++ {
//		mFit := make([]float64, len(rep.Pos))
//		for i := range rep.Pos {
//			mFit[i] = rep.PosFit[i][m]
//		}
//		idx := argsort(mFit)
//		mUp := make([]float64, len(mFit))
//		mDown := make([]float64, len(mFit))
//		for i := 0; i < len(mFit); i++ {
//			if i == 0 {
//				mDown[i] = math.Inf(1)
//			} else {
//				mDown[i] = mFit[idx[i-1]]
//			}
//			if i == len(mFit)-1 {
//				mUp[i] = math.Inf(1)
//			} else {
//				mUp[i] = mFit[idx[i+1]]
//			}
//		}
//		maxFit, minFit := mFit[idx[len(mFit)-1]], mFit[idx[0]]
//		for i := 0; i < len(mFit); i++ {
//			distance := (mUp[i] - mDown[i]) / (maxFit - minFit)
//			if math.IsNaN(distance) {
//				distance = math.Inf(1)
//			}
//			crowding[idx[i]] += distance
//		}
//	}
//	idx := argsort(crowding)
//	toDelete := idx[:nExtra]
//	keep := make([]bool, len(rep.Pos))
//	for i := range keep {
//		keep[i] = true
//	}
//	for _, i := range toDelete {
//		keep[i] = false
//	}
//	rep.Pos = filterRows(rep.Pos, keep, true)
//	rep.PosFit = filterRows(rep.PosFit, keep, true)
//	rep.PosPen = filterRows(rep.PosPen, keep, true)
//	updateGrid(rep, ngrid)
//}
//
//func mutation(POS [][]float64, gen, maxgen, Np int, varMax, varMin []float64, nVar int, u_mut float64) {
//	fract := float64(Np)/3.0 - math.Floor(float64(Np)/3.0)
//	var subSizes [3]int
//	if fract < 0.5 {
//		subSizes[0] = int(math.Ceil(float64(Np) / 3.0))
//		subSizes[1] = int(math.Round(float64(Np) / 3.0))
//		subSizes[2] = int(math.Round(float64(Np) / 3.0))
//	} else {
//		subSizes[0] = int(math.Round(float64(Np) / 3.0))
//		subSizes[1] = int(math.Round(float64(Np) / 3.0))
//		subSizes[2] = int(math.Floor(float64(Np) / 3.0))
//	}
//	cumSizes := [3]int{subSizes[0], subSizes[0] + subSizes[1], subSizes[0] + subSizes[1] + subSizes[2]}
//	// Uniform mutation
//	nmut := int(math.Round(u_mut * float64(subSizes[1])))
//	if nmut > 0 {
//		for i := 0; i < nmut; i++ {
//			idx := cumSizes[0] + rand.Intn(subSizes[1])
//			for j := 0; j < nVar; j++ {
//				POS[idx][j] = varMin[j] + rand.Float64()*(varMax[j]-varMin[j])
//			}
//		}
//	}
//	// Non-uniform mutation
//	perMut := math.Pow(1.0-float64(gen)/float64(maxgen), 5.0*float64(nVar))
//	nmut = int(math.Round(perMut * float64(subSizes[2])))
//	if nmut > 0 {
//		for i := 0; i < nmut; i++ {
//			idx := cumSizes[1] + rand.Intn(subSizes[2])
//			for j := 0; j < nVar; j++ {
//				POS[idx][j] = varMin[j] + rand.Float64()*(varMax[j]-varMin[j])
//			}
//		}
//	}
//}
//
//func checkBoundaries(POS, VEL [][]float64, maxVel, varMax, varMin []float64) {
//	Np := len(POS)
//	nVar := len(POS[0])
//	for i := 0; i < Np; i++ {
//		for j := 0; j < nVar; j++ {
//			if VEL[i][j] > maxVel[j] {
//				VEL[i][j] = maxVel[j]
//			}
//			if VEL[i][j] < -maxVel[j] {
//				VEL[i][j] = -maxVel[j]
//			}
//			if POS[i][j] > varMax[j] {
//				VEL[i][j] = -VEL[i][j]
//				POS[i][j] = varMax[j]
//			}
//			if POS[i][j] < varMin[j] {
//				VEL[i][j] = -VEL[i][j]
//				POS[i][j] = varMin[j]
//			}
//		}
//	}
//}
//
//// MOPSOAlgorithmNew is a new version matching the main interface
//// Archive is []*objectives.Result, agents are not stored
//// Use objectives.Problem for evaluation
//
//type MOPSOAlgorithmNew struct {
//	Params   MOPSOParams
//	Problem  objectives.Problem
//	Archive  []*objectives.Result
//}
//
//func NewMOPSOAlgorithmNew(problem objectives.Problem, params MOPSOParams) *MOPSOAlgorithmNew {
//	return &MOPSOAlgorithmNew{
//		Params:  params,
//		Problem: problem,
//	}
//}
//
//func (alg *MOPSOAlgorithmNew) Run() []*objectives.Result {
//	Np, Nr, MaxGen := alg.Params.Np, alg.Params.Nr, alg.Params.MaxGen
//	W, C1, C2 := alg.Params.W, alg.Params.C1, alg.Params.C2
//	ngrid, maxvel, u_mut := alg.Params.NGrid, alg.Params.MaxVel, alg.Params.UMut
//	varMin, varMax := alg.Problem.GetLowerBound(), alg.Problem.GetUpperBound()
//	nVar := alg.Problem.GetDimension()
//
//	// Initialization
//	POS := make([][]float64, Np)
//	VEL := make([][]float64, Np)
//	for i := 0; i < Np; i++ {
//		POS[i] = make([]float64, nVar)
//		VEL[i] = make([]float64, nVar)
//		for j := 0; j < nVar; j++ {
//			POS[i][j] = varMin[j] + rand.Float64()*(varMax[j]-varMin[j])
//		}
//	}
//	maxVelVec := make([]float64, nVar)
//	for i := 0; i < nVar; i++ {
//		maxVelVec[i] = (varMax[i] - varMin[i]) * maxvel / 100.0
//	}
//
//	// Evaluate initial population
//	POSFit := make([][]float64, Np)
//	results := make([]*objectives.Result, Np)
//	for i := 0; i < Np; i++ {
//		val, valuesWithKey, keys, _ := alg.Problem.Eval(POS[i])
//		POSFit[i] = val
//		results[i] = &objectives.Result{
//			Idx:      i,
//			Position: util.CopyArray(POS[i]),
//			Value:    util.CopyArray(val),
//			ValuesWithKey: valuesWithKey,
//			Key:      keys,
//		}
//	}
//
//	PBEST := deepCopyMatrix(POS)
//	PBESTFit := deepCopyMatrix(POSFit)
//
//	// Initial archive
//	dominated := checkDomination(POSFit)
//	archive := []*objectives.Result{}
//	for i := 0; i < Np; i++ {
//		if !dominated[i] {
//			archive = append(archive, results[i])
//		}
//	}
//	updateGridNew(archive, ngrid)
//
//	gen := 1
//	for gen <= MaxGen {
//		// Select leader
//		leader := selectLeaderNew(archive, ngrid)
//		leaderPos := leader.Position
//
//		// Update speeds and positions
//		for i := 0; i < Np; i++ {
//			for j := 0; j < nVar; j++ {
//				r1 := rand.Float64()
//				r2 := rand.Float64()
//				VEL[i][j] = W*VEL[i][j] +
//					C1*r1*(PBEST[i][j]-POS[i][j]) +
//					C2*r2*(leaderPos[j]-POS[i][j])
//				if math.Abs(VEL[i][j]) > maxVelVec[j] {
//					if VEL[i][j] > 0 {
//						VEL[i][j] = maxVelVec[j]
//					} else {
//						VEL[i][j] = -maxVelVec[j]
//					}
//				}
//				POS[i][j] += VEL[i][j]
//			}
//		}
//
//		// Mutation
//		mutationNew(POS, gen, MaxGen, Np, varMax, varMin, nVar, u_mut)
//
//		// Check boundaries
//		checkBoundaries(POS, VEL, maxVelVec, varMax, varMin)
//
//		// Evaluate
//		for i := 0; i < Np; i++ {
//			val, valuesWithKey, keys, _ := alg.Problem.Eval(POS[i])
//			POSFit[i] = val
//			results[i] = &objectives.Result{
//				Idx:      i,
//				Position: util.CopyArray(POS[i]),
//				Value:    util.CopyArray(val),
//				ValuesWithKey: valuesWithKey,
//				Key:      keys,
//			}
//		}
//
//		// Update archive
//		archive = updateArchiveNew(archive, results, ngrid)
//		if len(archive) > Nr {
//			archive = deleteFromArchiveNew(archive, len(archive)-Nr, ngrid)
//		}
//
//		// Update personal bests
//		for i := 0; i < Np; i++ {
//			if dominates(POSFit[i], PBESTFit[i]) {
//				PBESTFit[i] = util.CopyArray(POSFit[i])
//				PBEST[i] = util.CopyArray(POS[i])
//			} else if !dominates(PBESTFit[i], POSFit[i]) && rand.Float64() > 0.5 {
//				PBESTFit[i] = util.CopyArray(POSFit[i])
//				PBEST[i] = util.CopyArray(POS[i])
//			}
//		}
//		gen++
//	}
//	alg.Archive = archive
//	return archive
//}
//
//// --- Archive/grid/crowding helpers for new version ---
//
//func updateGridNew(archive []*objectives.Result, ngrid int) {
//	// No-op for now, can be extended for grid/crowding if needed
//}
//
//func selectLeaderNew(archive []*objectives.Result, ngrid int) *objectives.Result {
//	// Uniform random for now, can be extended for grid/crowding
//	return archive[rand.Intn(len(archive))]
//}
//
//func updateArchiveNew(archive []*objectives.Result, results []*objectives.Result, ngrid int) []*objectives.Result {
//	// Merge and keep only non-dominated
//	all := append(archive, results...)
//	dominated := checkDominationResults(all)
//	newArchive := []*objectives.Result{}
//	for i, res := range all {
//		if !dominated[i] {
//			newArchive = append(newArchive, res)
//		}
//	}
//	return newArchive
//}
//
//func deleteFromArchiveNew(archive []*objectives.Result, nExtra, ngrid int) []*objectives.Result {
//	// Use crowding distance truncation from objectives.DECD
//	return objectives.DECD(archive, nExtra)
//}
//
//func checkDominationResults(results []*objectives.Result) []bool {
//	Np := len(results)
//	dom := make([]bool, Np)
//	for i := 0; i < Np; i++ {
//		for j := 0; j < Np; j++ {
//			if i != j && dominates(results[j].Value, results[i].Value) {
//				dom[i] = true
//				break
//			}
//		}
//	}
//	return dom
//}
//
//// --- Deep copy helpers (replace util.CopyMatrix) ---
//func deepCopyMatrix(mat [][]float64) [][]float64 {
//	res := make([][]float64, len(mat))
//	for i := range mat {
//		res[i] = util.CopyArray(mat[i])
//	}
//	return res
//}
//
//// --- Argsort helper (replace util.Argsort) ---
//func argsort(x []float64) []int {
//	type kv struct{ k int; v float64 }
//	s := make([]kv, len(x))
//	for i, v := range x { s[i] = kv{i, v} }
//	sort.Slice(s, func(i, j int) bool { return s[i].v < s[j].v })
//	idx := make([]int, len(x))
//	for i, kv := range s { idx[i] = kv.k }
//	return idx
//}
