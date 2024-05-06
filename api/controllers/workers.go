package controller

import (
	"context"
	"routing/config"
	"routing/db"
)

func (c *Controller) getClosestNode(ctx context.Context, centroid interface{}, resultChan chan<- db.ClosestNodeResult) {
	defer close(resultChan)
	node, err := c.store.GetClosestPointToQueryLocation(ctx, centroid)
	select {
	case resultChan <- db.ClosestNodeResult{Node: node, Err: err}:
	case <-ctx.Done():
	}
}

func (c *Controller) getClosestNodeByUserLocationGeom(ctx context.Context, centroid interface{}, resultChan chan<- db.ClosestNodeToUserLocationResult) {
	defer close(resultChan)
	node, err := c.store.GetClosestPointToQueryLocationByLatLngGeom(ctx, centroid)
	select {
	case resultChan <- db.ClosestNodeToUserLocationResult{Node: node, Err: err}:
	case <-ctx.Done():
	}
}

func (c *Controller) calculateShortestPathWorker(ctx context.Context, fromID, toID int64, resultChan chan<- db.DijkstraResult) {
	defer close(resultChan)
	paths, distance, err := config.Dijkstra(c.Graph, fromID, toID)
	select {
	case resultChan <- db.DijkstraResult{Paths: paths, Distance: distance, Err: err}:
	case <-ctx.Done():
	}
}

func (c *Controller) getNodesByIdsWorker(ctx context.Context, ids []int64, resultChan chan<- db.Nodes) {
	defer close(resultChan)
	nodes, err := c.store.GetNodesByIds(ctx, ids)
	select {
	case resultChan <- db.Nodes{Nodes: nodes, Err: err}:
	case <-ctx.Done():
	}
}

func (c *Controller) getBuildingsWorker(ctx context.Context, buildingResult chan<- db.BuildingsResult) {
	defer close(buildingResult)
	building, err := c.store.ListBuildings(ctx)
	select {
	case buildingResult <- db.BuildingsResult{Buildings: building, Err: err}:
	case <-ctx.Done():
	}
}

func (c *Controller) getPlacesWorker(ctx context.Context, placesResult chan<- db.PlacesResult) {
	defer close(placesResult)
	places, err := c.store.ListPlaces(ctx)
	select {
	case placesResult <- db.PlacesResult{Places: places, Err: err}:
	case <-ctx.Done():
	}
}
