package api

import (
	"context"
	"routing/utils"
)

func (s *Server) ReadGraphIntoMemory(ctx context.Context) error {

	edges, err := s.store.ListEdges(ctx)
	if err != nil {
		return err
	}

	s.Graph = &utils.Graph{}
	for _, edge := range edges {
		err := s.Graph.AddEdge(edge.FromNodeID, edge.ToNodeID, edge.Weight)
		if err != nil {
			return err
		}
	}

	return err
}
