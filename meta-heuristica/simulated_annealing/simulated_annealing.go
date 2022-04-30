package simulated_annealing

import (
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot/plotter"
)

func (sai SimulatedAnnealingInstance) VerifySolution() bool {
	currentLoad := sai.ProblemInstance.GetInitialLoad()

	for _, step := range sai.Solution.WalkOrder {
		if currentLoad > sai.ProblemInstance.Limits[step] {
			return false
		}

		currentLoad -= sai.ProblemInstance.Demands[step]
	}

	return true
}

func (sai SimulatedAnnealingInstance) VerifySolution2(walkOrder []int) bool {
	currentLoad := sai.ProblemInstance.GetInitialLoad()

	for _, step := range walkOrder {
		if currentLoad > sai.ProblemInstance.Limits[step] {
			return false
		}

		currentLoad -= sai.ProblemInstance.Demands[step]
	}

	return true
}

func contains(n int, arr []int) bool {
	for _, el := range arr {
		if el == n {
			return true
		}
	}
	return false
}

func randIntExcept(min, max int, except []int) int {
	rand.Seed(time.Now().UnixNano())

	curr := rand.Intn(max-min) + min
	for contains(curr, except) {
		curr = rand.Intn(max-min) + min
	}

	return curr
}

func getArr(from, to int) []int {
	out := []int{}
	numb := to - from
	for i := 0; i < numb; i++ {
		out = append(out, from+i)
	}

	return out
}

func (sai SimulatedAnnealingInstance) GetFirstSolution() []int {

	nextPerm := func(p []int) {
		for i := len(p) - 1; i >= 0; i-- {
			if i == 0 || p[i] < len(p)-i-1 {
				p[i]++
				return
			}
			p[i] = 0
		}
	}

	getPerm := func(orig, p []int) []int {
		result := append([]int{}, orig...)
		for i, v := range p {
			result[i], result[i+v] = result[i+v], result[i]
		}
		return result
	}

	orig := getArr(1, sai.ProblemInstance.N)
	for p := make([]int, len(orig)); p[0] < len(p); nextPerm(p) {
		permutation := getPerm(orig, p)
		permutation = append(permutation, 0)
		if sai.VerifySolution2(permutation) {
			return permutation
		}
	}

	return nil
}

func (sai SimulatedAnnealingInstance) GenerateValidRandomInitialSolution(n int) []int {
	solution := generateInitialSolution(n)
	generatedSolutions := make([][]int, 0)
	sai.Solution.WalkOrder = solution
	generatedSolutions = append(generatedSolutions, solution)
	for true {
		solution = generateInitialSolution(n)
		if containsArr(solution, generatedSolutions) {
			continue
		}
		sai.Solution.WalkOrder = solution

		if sai.VerifySolution() {
			break
		}
	}

	return solution
}

func generateInitialSolution(n int) []int {
	sol := []int{}
	for i := 0; i < n-1; i++ {
		sol = append(sol, randIntExcept(1, n, sol))
	}

	sol = append(sol, 0)

	return sol
}

func (sai SimulatedAnnealingInstance) ObjectiveValue() int64 {
	totalDistance := sai.ProblemInstance.Distances[0][sai.Solution.WalkOrder[0]]
	for i := 0; i < sai.ProblemInstance.N-1; i++ {
		a := sai.Solution.WalkOrder[i]
		b := sai.Solution.WalkOrder[i+1]
		distance := sai.ProblemInstance.Distances[a][b]
		totalDistance += distance
	}

	return int64(totalDistance)
}

func (sai SimulatedAnnealingInstance) ObjectiveValueWithSolution(walkOrder []int) int64 {
	totalDistance := sai.ProblemInstance.Distances[0][walkOrder[0]]
	for i := 0; i < sai.ProblemInstance.N-1; i++ {
		a := walkOrder[i]
		b := walkOrder[i+1]
		distance := sai.ProblemInstance.Distances[a][b]
		totalDistance += distance
	}

	return int64(totalDistance)
}

func isEqual(arr1, arr2 []int) bool {
	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func containsArr(arr []int, mat [][]int) bool {
	for i := 0; i < len(mat); i++ {
		if isEqual(arr, mat[i]) {
			return true
		}
	}

	return false
}

func (sai SimulatedAnnealingInstance) GetNeighbor(except [][]int) []int {
	neighbor := make([]int, len(sai.Solution.WalkOrder))
	for i := 0; i < len(sai.Solution.WalkOrder); i++ {
		for j := 0; j < len(sai.Solution.WalkOrder); j++ {
			copied := make([]int, len(neighbor))

			copy(copied, sai.Solution.WalkOrder)

			a := copied[i]
			b := copied[j]

			copied[i] = b
			copied[j] = a

			if isEqual(copied, sai.Solution.WalkOrder) || containsArr(copied, except) {
				continue
			}

			valid := sai.VerifySolution2(copied)
			if valid {
				return copied
			}
		}
	}

	return sai.Solution.WalkOrder
}

func (sai SimulatedAnnealingInstance) GetNeighbors() [][]int {
	neighbors := make([][]int, 0)
	for i := 0; i < len(sai.Solution.WalkOrder); i++ {
		for j := 0; j < len(sai.Solution.WalkOrder); j++ {
			copied := make([]int, len(sai.Solution.WalkOrder))

			copy(copied, sai.Solution.WalkOrder)

			a := copied[i]
			b := copied[j]

			if a == 0 || b == 0 {
				continue
			}

			copied[i] = b
			copied[j] = a

			if isEqual(copied, sai.Solution.WalkOrder) {
				continue
			}

			valid := sai.VerifySolution2(copied)
			if valid {
				neighbors = append(neighbors, copied)
			}
		}
	}

	return neighbors
}

func randInt(min, max int) int {
	max++
	if max-min == 0 {
		return 0
	}
	return rand.Intn(max-min) + min
}

func getRandomNeighborExcept(neighbors [][]int, except [][]int) []int {
	toGet := randInt(0, len(neighbors)-1)
	for containsArr(neighbors[toGet], except) {
		toGet = randInt(0, len(neighbors))
	}

	return neighbors[toGet]
}

func (sai *SimulatedAnnealingInstance) Solve(timeLimit int64, params SimulatedAnnealingParams) {
	sai.Solution = &Solution{
		BKV:       0,
		WalkOrder: make([]int, 0),
	}
	sai.Solution.WalkOrder = sai.GenerateValidRandomInitialSolution(sai.ProblemInstance.N)
	sai.Solution.BKV = sai.ObjectiveValue()

	T := float32(params.T0)
	j := 1
	successes := 0

	solutions := [][]int{}
	for true {
		i := 1
		successes = 0

		neighbors := sai.GetNeighbors()
		iterationSolutions := make([][]int, 0)
		for true {
			used := [][]int{}
			si := getRandomNeighborExcept(neighbors, used)
			siVal := sai.ObjectiveValueWithSolution(si)
			currVal := sai.ObjectiveValue()
			delta := siVal - currVal
			if delta <= 0 || (math.Exp(float64(float32(-delta)/(T)))) > float64(randInt(0, 1)) {
				// accepts new solution
				sai.Solution.WalkOrder = si
				sai.Solution.BKV = siVal
				successes++
				iterationSolutions = append(iterationSolutions, si)
			}

			i += 1
			used = append(used, si)
			if successes >= params.L || i > params.P || len(used) >= sai.ProblemInstance.N*sai.ProblemInstance.N /* should stop if all neighbors were tested already */ {
				break
			}
		}

		sai.Points = append(sai.Points, plotter.XY{
			X: float64(T),
			Y: float64(sai.Solution.BKV),
		})

		T = params.Alpha * T
		j += 1

		solutions = append(solutions, sai.Solution.WalkOrder)

		if successes == 0 || T < float32(params.TF) {
			break
		}
	}

	sai.Solution.WalkOrder = sai.getBest(solutions)

	sai.Solved = true
	sai.Solution.BKV = sai.ObjectiveValue()
}

func (sai SimulatedAnnealingInstance) getBest(solutions [][]int) []int {
	if len(solutions) == 0 {
		return sai.Solution.WalkOrder
	}
	best := 0
	distance := sai.ObjectiveValueWithSolution(solutions[0])

	for index, solution := range solutions {
		if index == 0 {
			continue
		}
		currDistance := sai.ObjectiveValueWithSolution(solution)
		if currDistance < int64(distance) {
			distance = currDistance
			best = index
		}
	}

	return solutions[best]
}
