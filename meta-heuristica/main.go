package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/reonardoleis/simulated-annealing-tsp-mod/models"
	"github.com/reonardoleis/simulated-annealing-tsp-mod/simulated_annealing"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
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
		T0:              40,    // Defines the initial temperature
		TF:              0.001, // Defines the final temperature
		P:               125,   // Defines number of metropolis calls per cycle
		MaxAcceptances:  50,    // Defines how the maximum of (maybe) better solutions that can be found each metropolis call
		MetropolisLimit: 175,   // Sets a limit for the number of iterations on metropolis method
		Alpha:           0.99,  // Defines the temperature "decay" coefficient
	})

	fmt.Println("Smallest route: ", simulatedAnnealingInstance.Solution.TraveledDistance)

	// Generates the plot
	p := plot.New()
	p.Title.Text = "Objective Value x Temperature\nRed line: Iteration accepted\nBlue line: Metropolis accepted"
	p.X.Label.Text = "Temperature"
	p.Y.Label.Text = "Objective Value"
	p.Add(plotter.NewGrid())

	l, err := plotter.NewLine(simulatedAnnealingInstance.MetropolisPoints)
	if err != nil {
		log.Fatalln(err.Error())
	}

	l.LineStyle.Color = color.RGBA{R: 0, G: 0, B: 255, A: 255}
	p.Add(l)

	l, err = plotter.NewLine(*&simulatedAnnealingInstance.Points)
	if err != nil {
		log.Fatalln(err.Error())
	}

	l.LineStyle.Width = l.LineStyle.Width * 2
	l.LineStyle.Color = color.RGBA{R: 255, G: 0, B: 0, A: 255}
	p.Add(l)

	if err = p.Save(10*vg.Inch, 4*vg.Inch, "plot.png"); err != nil {
		log.Fatalln(err.Error())
	}

}
