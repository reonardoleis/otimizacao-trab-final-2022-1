package simulated_annealing

import (
	"github.com/reonardoleis/simulated-annealing-tsp-mod/models"
	"gonum.org/v1/plot/plotter"
)

type SimulatedAnnealingInstance struct {
	ProblemInstance models.Instance
	Solution        *Solution
	Solved          bool
	Points          plotter.XYs
}

type Solution struct {
	BKV       int64
	WalkOrder []int
}

type SimulatedAnnealingParams struct {
	T0    int
	TF    float32
	P     int
	L     int
	Alpha float32
}

func NewSimulatedAnnealingInstance(problemInstance models.Instance) SimulatedAnnealingInstance {
	return SimulatedAnnealingInstance{
		ProblemInstance: problemInstance,
		Solution: &Solution{
			BKV:       0,
			WalkOrder: []int{},
		},
		Solved: false,
		Points: make(plotter.XYs, 0),
	}
}
