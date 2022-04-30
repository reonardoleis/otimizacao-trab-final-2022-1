package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/reonardoleis/simulated-annealing-tsp-mod/models"
	"github.com/reonardoleis/simulated-annealing-tsp-mod/simulated_annealing"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	instance := models.Instance{}
	err := instance.LoadFromFile("./instances/instance_21.dat")
	if err != nil {
		log.Fatalln(err.Error())
	}

	simulatedAnnealingInstance := simulated_annealing.NewSimulatedAnnealingInstance(instance)

	simulatedAnnealingInstance.Solve(0, simulated_annealing.SimulatedAnnealingParams{
		T0:    30,    // Defines the initial temperature
		TF:    0.001, // Defines the final temperature
		P:     1000,  // Defines number of iterations per cycle, will be upper bounded by the maximum number of neighbors
		L:     500,   // Defines how the maximum of (maybe) better solutions that can be found each iteration
		Alpha: 0.99,  // Defines the temperature "decay" coefficient
	})

	fmt.Println("Smallest route: ", simulatedAnnealingInstance.Solution.BKV)

	// Generates the plot
	p := plot.New()
	p.Title.Text = "Objective Value x Temperature"
	p.X.Label.Text = "Temperature"
	p.Y.Label.Text = "Objective Value"
	p.Add(plotter.NewGrid())
	err = plotutil.AddLines(p, "", simulatedAnnealingInstance.Points)
	if err != nil {
		log.Fatalln(err.Error())
	}

	if err = p.Save(10*vg.Inch, 4*vg.Inch, "plot.png"); err != nil {
		log.Fatalln(err.Error())
	}

}
