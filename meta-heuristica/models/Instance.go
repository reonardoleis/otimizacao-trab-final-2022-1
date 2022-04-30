package models

import (
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"
)

type Instance struct {
	N         int     // Vertices count
	E         int     // Edges count, complete graph
	Distances [][]int // Distances matrix of the graph
	Demands   []int   // Demand of each vertex
	Limits    []int   // Load limit of each vertex
}

func (inst *Instance) LoadFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	content := strings.Split(strings.ReplaceAll(string(b), "\r", ""), "\n")

	inst.Distances = make([][]int, 0)
	inst.Demands = make([]int, 0)
	inst.Limits = make([]int, 0)

	n, err := strconv.Atoi(content[0])
	if err != nil {
		return err
	}

	inst.N = n
	inst.E = int(math.Floor(float64(inst.N * (inst.N - 1) / 2)))
	for i := 1; i != inst.N+1; i++ {
		inst.Distances = append(inst.Distances, []int{})
		lineDistances := strings.Split(content[i], " ")
		for _, lineDistance := range lineDistances {
			if lineDistance == "" {
				continue
			}
			n, err := strconv.Atoi(lineDistance)
			if err != nil {
				return err
			}
			inst.Distances[i-1] = append(inst.Distances[i-1], n)
		}
	}

	demands := strings.Split(content[len(content)-3], " ")
	limits := strings.Split(content[len(content)-2], " ")

	for index, _ := range demands {
		if demands[index] != "" {
			demand, err := strconv.Atoi(demands[index])
			if err != nil {
				return err
			}
			inst.Demands = append(inst.Demands, demand)
		}

		if limits[index] != "" {
			limit, err := strconv.Atoi(limits[index])
			if err != nil {
				return err
			}

			inst.Limits = append(inst.Limits, limit)
		}
	}

	return nil
}

func (inst Instance) GetInitialLoad() int {
	sum := 0
	for _, demand := range inst.Demands {
		sum += demand
	}

	return sum
}
