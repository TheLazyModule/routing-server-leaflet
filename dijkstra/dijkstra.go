package dijkstra

import (
	"fmt"
	"math"
	rg "routing/graph"
)

type PreviousNodeAndWeight struct {
	prevNode string
	weight   float64
}

// Dijkstra finds the shortest path from initial to end node and returns the path.
func Dijkstra(graph *rg.Graph, initial, end string) ([]string, error) {
	shortestPaths := make(map[string]PreviousNodeAndWeight)
	shortestPaths[initial] = PreviousNodeAndWeight{prevNode: "", weight: 0}
	visited := make(map[string]bool)
	currentNode := initial

	for currentNode != "" {
		if currentNode == end {
			break
		}
		visited[currentNode] = true
		destinations := graph.GetEdges()[currentNode]
		weightToCurrentNode := shortestPaths[currentNode].weight

		for _, nextNode := range destinations {
			weight := graph.GetWeights()[rg.NodePair{currentNode, nextNode}] + weightToCurrentNode
			if nextWeight, ok := shortestPaths[nextNode]; !ok || weight < nextWeight.weight {
				shortestPaths[nextNode] = PreviousNodeAndWeight{prevNode: currentNode, weight: weight}
			}
		}

		currentNode = ""
		minWeight := math.Inf(1)
		for node, pw := range shortestPaths {
			if !visited[node] && pw.weight < minWeight {
				minWeight = pw.weight
				currentNode = node
			}
		}
		if currentNode == "" {
			return nil, fmt.Errorf("route not possible")
		}
	}

	// Reconstruct path
	var path []string
	for curr := end; curr != ""; curr = shortestPaths[curr].prevNode {
		path = append([]string{curr}, path...)
	}

	return path, nil
}
