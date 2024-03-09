package api

import (
	"context"
	"routing/utils"
)

func (s *Server) ReadGraphIntoMemory(ctx context.Context) (*utils.Graph, error) {

	edges, err := s.store.ListEdges(ctx)
	if err != nil {
		return nil, err
	}

	weights, err := s.store.ListWeights(ctx)
	if err != nil {
		return nil, err
	}

	newGraph := utils.NewGraph()
	err = utils.ReadIntoMemory(newGraph, edges)
	if err != nil {
		return nil, err
	}

	err = utils.ReadIntoMemory(newGraph, weights)
	if err != nil {
		return nil, err
	}
	return newGraph, err
}
