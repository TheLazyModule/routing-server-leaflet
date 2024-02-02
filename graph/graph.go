package graph

/*
	Graph

edges is a map of all possible next nodes

	e.g. {'X': ['A', 'B', 'C', 'E'], ...}
	weights has all the weights between two nodes,
	with the two nodes as a tuple as the key
	e.g. {[ 'X', 'A' ]: 7, [ 'X', 'B' ]: 2, ...}
*/

type NodePair [2]string

type Node struct {
	label string
	x, y  float64
}

type Graph struct {
	edges   map[string][]string
	weights map[NodePair]float64
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) AddEdge(fromNode, toNode string, weight float64) {
	if g.edges == nil {
		g.edges = make(map[string][]string)
	}
	if g.weights == nil {
		g.weights = make(map[NodePair]float64)
	}
	g.edges[fromNode] = append(g.edges[fromNode], toNode)
	g.edges[toNode] = append(g.edges[toNode], fromNode)
	g.weights[NodePair{fromNode, toNode}] = weight
}

func (g *Graph) GetEdges() map[string][]string {
	return g.edges
}

func (g *Graph) GetWeights() map[NodePair]float64 {
	return g.weights
}
