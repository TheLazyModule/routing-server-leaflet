package main

import (
	"fmt"
	"routing/dijkstra"
	rg "routing/graph"
)

func main() {
	//var newGraph = dk.Graph[int]{}
	var newGraph = rg.NewGraph()

	type edges struct {
		fromNode, toNode string
		weight           float64
	}

	var edgesContainer = []edges{
		{"X", "A", 7},
		{"A", "B", 3},
		{"A", "D", 4},
		{"X", "B", 2},
		{"X", "C", 3},
		{"B", "H", 5},
		{"B", "D", 4},
		{"C", "L", 2},
		{"D", "F", 1},
		{"X", "E", 4},
		{"F", "H", 3},
		{"G", "Y", 2},
		{"G", "H", 2},
		{"I", "J", 6},
		{"I", "L", 4},
		{"I", "K", 4},
		{"K", "Y", 5},
		{"J", "L", 1},
	}

	for _, edge := range edgesContainer {
		newGraph.AddEdge(edge.fromNode, edge.toNode, edge.weight)
	}

	//fmt.Println(newGraph.GetWeights())
	//fmt.Println(newGraph.GetEdges())
	fmt.Println(dijkstra.Dijkstra(newGraph, "X", "Y"))

}
