package utils

import (
	"fmt"
	"math"
)

type PreviousNodeAndWeight struct {
	prevNode int64 // Adjusted to int64 to match node IDs
	weight   float64
}

// Dijkstra finds the shortest path from initial to end node and returns the path as node IDs.
func Dijkstra(graph *Graph, initial, end int64) ([]int64, float64, error) {
	shortestPaths := make(map[int64]PreviousNodeAndWeight)
	shortestPaths[initial] = PreviousNodeAndWeight{prevNode: -1, weight: 0} // Use -1 to indicate no previous node
	visited := make(map[int64]bool)
	currentNode := initial

	for currentNode != -1 {
		if currentNode == end {
			break
		}
		visited[currentNode] = true
		destinations := graph.GetEdges()[currentNode]
		weightToCurrentNode := shortestPaths[currentNode].weight

		for _, nextNode := range destinations {
			weight := graph.GetWeights()[NodePair{currentNode, nextNode}] + weightToCurrentNode
			if nextWeight, ok := shortestPaths[nextNode]; !ok || weight < nextWeight.weight {
				shortestPaths[nextNode] = PreviousNodeAndWeight{prevNode: currentNode, weight: weight}
			}
		}

		currentNode = -1
		minWeight := math.Inf(1)
		for node, pw := range shortestPaths {
			if !visited[node] && pw.weight < minWeight {
				minWeight = pw.weight
				currentNode = node
			}
		}
	}

	if shortestPaths[end].prevNode == -1 {
		return nil, 0, fmt.Errorf("route not possible")
	}

	// Reconstruct path
	var path []int64
	for curr := end; curr != -1; curr = shortestPaths[curr].prevNode {
		path = append([]int64{curr}, path...)
	}

	return path, shortestPaths[end].weight, nil
}
