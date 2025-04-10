package gwo

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"golang-moaha-construction/internal/algorithms"
	"golang-moaha-construction/internal/data"
	"golang-moaha-construction/internal/objectives"
	"math"
	"math/rand"
	"sync"
)

const NameType algorithms.AlgorithmType = "GWO"

const (
	numberOfObjective = 1
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

type Config struct {
	NumberOfAgents int
	NumberOfIter   int
	AParam         float64
}

func Create(
	problem objectives.Problem,
	configs Config,
) (*GWOAlgorithm, error) {

	convergence := make([]float64, configs.NumberOfIter)
	agents := make([]*objectives.Result, configs.NumberOfAgents)

	if numberOfObjective != problem.NumberOfObjectives() {
		return nil, objectives.ErrInvalidNumberOfObjectives
	}

	return &GWOAlgorithm{
		NumberOfAgents:    configs.NumberOfAgents,
		NumberOfIter:      configs.NumberOfIter,
		AParam:            configs.AParam,
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
				value, _, _ := g.ObjectiveFunction.Eval(g.Agents[agentIdx].Position)

				g.Agents[agentIdx].Value = value
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

func (g *GWOAlgorithm) RunWithChannel(doneChan chan<- struct{}, channel chan<- any) error {

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
				value, _, _ := g.ObjectiveFunction.Eval(g.Agents[agentIdx].Position)

				g.Agents[agentIdx].Value = value
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
			newAgent := &objectives.Result{
				Idx:      agentIdx,
				Position: positions,
			}

			value, _, _ := g.ObjectiveFunction.Eval(positions)
			newAgent.Value = value

			g.Agents[agentIdx] = newAgent
		}(agentIdx)
	}
	wg.Wait()

	g.findBest()
}

// This algorithms solve only one objective
func (g *GWOAlgorithm) findBest() {
	for i := range g.Agents {
		if g.Agents[i].Value[0] < g.Alpha.Value[0] {
			g.Alpha = g.Agents[i].CopyAgent()
		} else if g.Agents[i].Value[0] < g.Beta.Value[0] {
			g.Beta = g.Agents[i].CopyAgent()
		} else if g.Agents[i].Value[0] < g.Gamma.Value[0] {
			g.Gamma = g.Agents[i].CopyAgent()
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
