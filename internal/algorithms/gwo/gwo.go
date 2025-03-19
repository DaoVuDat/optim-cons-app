package gwo

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"golang-moaha-construction/internal/util"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
)

const (
	numberOfObjective = 1
	Name              = "GWO"
	NUM_AGENTS        = "Number of Agents"
	NUM_ITERS         = "Number of Iteration"
	PARAM_A           = "Parameter a"
)

type GWOAlgorithm struct {
	NumberOfAgents    int
	NumberOfIter      int
	Agents            []*objectives.Result
	AParam            float64
	Alpha             *objectives.Result
	Beta              *objectives.Result
	Gamma             *objectives.Result
	Convergence       []float64
	ObjectiveFunction objectives.Problem
}

var Configs = []data.Config{
	{
		Name:               NUM_AGENTS,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               NUM_ITERS,
		ValidationFunction: util.IsValidPositiveInteger,
	},
	{
		Name:               PARAM_A,
		ValidationFunction: util.IsValidFloat,
	},
}

func Create(
	problem objectives.Problem,
	configs []*data.Config,
) (*GWOAlgorithm, error) {

	var numsIters, numAgents int
	var aParams float64

	for _, config := range configs {
		// sanity check
		val := strings.Trim(config.Value, " ")
		switch config.Name {
		case NUM_AGENTS:
			numAgents, _ = strconv.Atoi(val)
		case NUM_ITERS:
			numsIters, _ = strconv.Atoi(val)
		case PARAM_A:
			aParams, _ = strconv.ParseFloat(val, 64)
		}
	}

	convergence := make([]float64, numsIters)
	agents := make([]*objectives.Result, numAgents)

	if numberOfObjective != problem.NumberOfObjectives() {
		return nil, objectives.ErrInvalidNumberOfObjectives
	}

	return &GWOAlgorithm{
		NumberOfAgents:    numAgents,
		NumberOfIter:      numsIters,
		AParam:            aParams,
		Convergence:       convergence,
		ObjectiveFunction: problem,
		Agents:            agents,
	}, nil
}

func (g *GWOAlgorithm) Type() data.TypeProblem {
	return data.Single
}

func (g *GWOAlgorithm) Run() error {

	// initialization
	g.initialization()

	l := 0
	a := g.AParam

	bar := progressbar.Default(int64(g.NumberOfIter))
	var wg sync.WaitGroup

	for l < g.NumberOfIter {
		a = 2.0 - float64(l)*(2.0/float64(g.NumberOfIter))

		wg.Add(g.NumberOfAgents)
		for agentIdx := range g.Agents {
			go func(agentIdx int) {
				defer wg.Done()

				for posIdx := range g.Agents[agentIdx].Position {
					// Alpha
					r1 := rand.Float64()
					r2 := rand.Float64()
					A := 2*a*r1 - a
					C := 2 * r2
					D := math.Abs(C*g.Alpha.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
					XAlpha := g.Alpha.Position[posIdx] - A*D

					// Beta
					r1 = rand.Float64()
					r2 = rand.Float64()
					A = 2*a*r1 - a
					C = 2 * r2
					D = math.Abs(C*g.Beta.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
					XBeta := g.Beta.Position[posIdx] - A*D

					// Gamma
					r1 = rand.Float64()
					r2 = rand.Float64()
					A = 2*a*r1 - a
					C = 2 * r2
					D = math.Abs(C*g.Gamma.Position[posIdx] - g.Agents[agentIdx].Position[posIdx])
					XGamma := g.Gamma.Position[posIdx] - A*D

					g.Agents[agentIdx].Position[posIdx] = (XAlpha + XBeta + XGamma) / 3
				}
				// check out of boundaries
				g.outOfBoundaries(g.Agents[agentIdx].Position)

				// evaluate
				g.Agents[agentIdx] = g.ObjectiveFunction.Eval(g.Agents[agentIdx].Position)

			}(agentIdx)
		}

		wg.Wait()

		g.findBest()

		g.Convergence[l] = g.Alpha.Value[0]
		bar.Describe(fmt.Sprintf("Iter %d: %e", l+1, g.Alpha.Value[0]))
		bar.Add(1)

		l++
	}

	return nil
}

func (g *GWOAlgorithm) initialization() {

	vals := make([]float64, g.ObjectiveFunction.NumberOfObjectives())
	for i := 0; i < g.ObjectiveFunction.NumberOfObjectives(); i++ {
		if g.ObjectiveFunction.FindMin() {
			vals[i] = math.MaxFloat64
		} else {
			vals[i] = math.MinInt64
		}
	}

	g.Alpha = &objectives.Result{
		Value: vals,
	}

	g.Beta = &objectives.Result{
		Value: vals,
	}

	g.Gamma = &objectives.Result{
		Value: vals,
	}

	var wg sync.WaitGroup
	wg.Add(g.NumberOfAgents)
	for agentIdx := range g.Agents {
		go func(agentIdx int) {
			defer wg.Done()
			positions := make([]float64, g.ObjectiveFunction.GetDimension())
			for i := 0; i < g.ObjectiveFunction.GetDimension(); i++ {
				positions[i] = g.ObjectiveFunction.GetLowerBound()[i] + rand.Float64()*
					(g.ObjectiveFunction.GetUpperBound()[i]-g.ObjectiveFunction.GetLowerBound()[i])
			}

			// evaluate
			g.Agents[agentIdx] = g.ObjectiveFunction.Eval(positions)
			g.Agents[agentIdx].Idx = agentIdx
		}(agentIdx)
	}
	wg.Wait()

	g.findBest()
}

// This algorithm solve only one objective
func (g *GWOAlgorithm) findBest() {
	for i := range g.Agents {
		if g.Agents[i].Value[0] < g.Alpha.Value[0] {
			g.Alpha = util.CopyAgent(g.Agents[i])
		} else if g.Agents[i].Value[0] < g.Beta.Value[0] {
			g.Beta = util.CopyAgent(g.Agents[i])
		} else if g.Agents[i].Value[0] < g.Gamma.Value[0] {
			g.Gamma = util.CopyAgent(g.Agents[i])
		}
	}
}

func (g *GWOAlgorithm) outOfBoundaries(pos []float64) {
	for i := range pos {
		if pos[i] < g.ObjectiveFunction.GetLowerBound()[i] {
			pos[i] = g.ObjectiveFunction.GetLowerBound()[i]
		} else if pos[i] > g.ObjectiveFunction.GetUpperBound()[i] {
			pos[i] = g.ObjectiveFunction.GetUpperBound()[i]
		}
	}
}
