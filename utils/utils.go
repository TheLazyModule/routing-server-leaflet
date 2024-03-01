package utils

import (
	"fmt"
	"routing/db/sqlc"
)

func ReadIntoMemory(graph *Graph, edges interface{}) error {
	// Type switch to handle different types
	switch e := edges.(type) {
	case []db.ListEdgesRow:
		for _, edge := range e {
			if err := graph.AddEdgesFromDB(edge.NodeID, edge.Neighbors); err != nil {
				return err
			}
		}
	case []db.Weight:
		for _, weight := range e {
			if err := graph.AddWeightsFromDB(weight.FromNodeID, weight.ToNodeID, weight.Distance); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unsupported type passed to ReadIntoMemory")
	}

	return nil
}
