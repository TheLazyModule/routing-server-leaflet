package utils

import "routing/db/dto"

/*
	Graph

edges is a map of all possible next nodes

	e.g. {'X': ['A', 'B', 'C', 'E'], ...}
	weights has all the weights between two nodes,
	with the two nodes as a tuple as the key
	e.g. {[ 'X', 'A' ]: 7, [ 'X', 'B' ]: 2, ...}
*/

type NodePair [2]int64

type Node struct {
	label string
	x, y  float64
}

type Graph struct {
	edges   map[int64]dto.EdgesData
	weights map[NodePair]float64
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g *Graph) AddEdge(fromNode, toNode int64, weight float64) error {
	g.initializeMaps()
	g.edges[fromNode] = append(g.edges[fromNode], toNode)
	g.edges[toNode] = append(g.edges[toNode], fromNode)
	g.weights[NodePair{fromNode, toNode}] = weight
	return nil
}

func (g *Graph) initializeMaps() {
	if g.edges == nil {
		g.edges = make(map[int64]dto.EdgesData)
	}
	if g.weights == nil {
		g.weights = make(map[NodePair]float64)
	}
}

func (g *Graph) GetEdges() map[int64]dto.EdgesData {
	return g.edges
}

func (g *Graph) GetWeights() map[NodePair]float64 {
	return g.weights
}

// AddEdgesFromDB adds edges from db
func (g *Graph) AddEdgesFromDB(nodeId int64, neighbors dto.EdgesData) error {
	g.initializeMaps()
	g.edges[nodeId] = neighbors
	return nil
}

// AddWeightsFromDB adds weights from db
func (g *Graph) AddWeightsFromDB(fromNodeID, toNodeID int64, distance float64) error {
	g.initializeMaps()
	g.weights[NodePair{fromNodeID, toNodeID}] = distance
	return nil
}
