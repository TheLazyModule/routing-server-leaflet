package api

import (
	"context"
	"routing/db"
	"routing/utils"
)

func (s *Server) getClosestNode(ctx context.Context, centroid interface{}, resultChan chan<- db.ClosestNodeResult) {
	defer close(resultChan)
	node, err := s.store.GetClosestPointToQueryLocation(ctx, centroid)
	select {
	case resultChan <- db.ClosestNodeResult{Node: node, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) getClosestNodeByUserLocationGeom(ctx context.Context, centroid interface{}, resultChan chan<- db.ClosestNodeToUserLocationResult) {
	defer close(resultChan)
	node, err := s.store.GetClosestPointToQueryLocationByLatLngGeom(ctx, centroid)
	select {
	case resultChan <- db.ClosestNodeToUserLocationResult{Node: node, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) calculateShortestPathWorker(ctx context.Context, fromID, toID int64, resultChan chan<- db.DijkstraResult) {
	defer close(resultChan)
	paths, distance, err := utils.Dijkstra(s.Graph, fromID, toID)
	select {
	case resultChan <- db.DijkstraResult{Paths: paths, Distance: distance, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) getNodesByIdsWorker(ctx context.Context, ids []int64, resultChan chan<- db.Nodes) {
	defer close(resultChan)
	nodes, err := s.store.GetNodesByIds(ctx, ids)
	select {
	case resultChan <- db.Nodes{Nodes: nodes, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) getBuildingsWorker(ctx context.Context, buildingResult chan<- db.BuildingsResult) {
	defer close(buildingResult)
	building, err := s.store.ListBuildings(ctx)
	select {
	case buildingResult <- db.BuildingsResult{Buildings: building, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) getPlacesWorker(ctx context.Context, placesResult chan<- db.PlacesResult) {
	defer close(placesResult)
	places, err := s.store.ListPlaces(ctx)
	select {
	case placesResult <- db.PlacesResult{Places: places, Err: err}:
	case <-ctx.Done():
	}
}
