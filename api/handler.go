package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/db"
	"time"
)

const UserLocation = "My Location"

func (s *Server) ShowMap(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/map")
}

func (s *Server) GetPlacesAndBuildings(ctx *gin.Context) {
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get places and buildings asynchronously.
	placesChan := make(chan db.PlacesResult, 1)
	// -->
	go s.getPlacesWorker(pipelineCtx, placesChan)

	buildingsChan := make(chan db.BuildingsResult, 1)
	// -->
	go s.getBuildingsWorker(pipelineCtx, buildingsChan)

	placesResult := <-placesChan
	if err := placesResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	buildingsResult := <-buildingsChan
	if err := buildingsResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"places": placesResult.Places, "buildings": buildingsResult.Buildings})
}

func (s *Server) GetShortestRouteByBuildingOrPlace(ctx *gin.Context) {
	var req db.RouteRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeFromUserLocationChan := make(chan db.ClosestNodeToUserLocationResult, 1)

	if req.From.String == UserLocation {
		// -->
		go s.getClosestNodeByUserLocationGeom(pipelineCtx, req.FromLocation, closestNodeFromUserLocationChan)
	} else {

		geomFrom, err := s.store.GetBuildingOrPlace(pipelineCtx, req.From.String)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		// -->
		go s.getClosestNode(pipelineCtx, geomFrom.Geom, closestNodeFromChan)
	}

	// Get Building or Place for GeomTo
	geomTo, err := s.store.GetBuildingOrPlace(pipelineCtx, req.To.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	// -->
	go s.getClosestNode(pipelineCtx, geomTo.Geom, closestNodeToChan)

	timeout := 10 * time.Second
	select {
	case <-time.After(timeout):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
		return
	case closestNodeFromUserLocationResult := <-closestNodeFromUserLocationChan:
		if closestNodeFromUserLocationResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		closestNodeToResult := <-closestNodeToChan

		if err = closestNodeToResult.Err; err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		dijkstraResultChan := make(chan db.DijkstraResult, 1)
		// -->
		go s.calculateShortestPathWorker(pipelineCtx, closestNodeFromUserLocationResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

		dijkstraResult := <-dijkstraResultChan
		if dijkstraResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
			return
		}

		nodesChan := make(chan db.Nodes, 1)
		// -->
		go s.getNodesByIdsWorker(pipelineCtx, dijkstraResult.Paths, nodesChan)

		nodesResult := <-nodesChan
		if nodesResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
		return
	case closestNodeFromResult := <-closestNodeFromChan:
		if closestNodeFromResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		closestNodeToResult := <-closestNodeToChan

		if err = closestNodeToResult.Err; err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		dijkstraResultChan := make(chan db.DijkstraResult, 1)
		// -->
		go s.calculateShortestPathWorker(pipelineCtx, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

		dijkstraResult := <-dijkstraResultChan
		if dijkstraResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
			return
		}

		nodesChan := make(chan db.Nodes, 1)
		// -->
		go s.getNodesByIdsWorker(pipelineCtx, dijkstraResult.Paths, nodesChan)

		nodesResult := <-nodesChan
		if nodesResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
		return
	}
}
