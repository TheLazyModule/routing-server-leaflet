package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/api/utils"
	"routing/db"
	"time"
)

const UserLocation = "My Location"

func (c *Controller) ShowMap(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/map")
}

func (c *Controller) GetPlacesAndBuildings(ctx *gin.Context) {
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get places and buildings asynchronously.
	placesChan := make(chan db.PlacesResult, 1)
	// -->
	go c.getPlacesWorker(pipelineCtx, placesChan)

	buildingsChan := make(chan db.BuildingsResult, 1)
	// -->
	go c.getBuildingsWorker(pipelineCtx, buildingsChan)

	placesResult := <-placesChan
	if err := placesResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	buildingsResult := <-buildingsChan
	if err := buildingsResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"places": placesResult.Places, "buildings": buildingsResult.Buildings})
}

func (c *Controller) GetShortestRouteByBuildingOrPlace(ctx *gin.Context) {
	var req db.RouteRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeFromUserLocationChan := make(chan db.ClosestNodeToUserLocationResult, 1)

	if req.From.String == UserLocation {
		// -->
		go c.getClosestNodeByUserLocationGeom(pipelineCtx, req.FromLocation, closestNodeFromUserLocationChan)
	} else {

		geomFrom, err := c.store.GetBuildingOrPlace(pipelineCtx, req.From.String)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
		// -->
		go c.getClosestNode(pipelineCtx, geomFrom.Geom, closestNodeFromChan)
	}

	// Get Building or Place for GeomTo
	geomTo, err := c.store.GetBuildingOrPlace(pipelineCtx, req.To.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	// -->
	go c.getClosestNode(pipelineCtx, geomTo.Geom, closestNodeToChan)

	timeout := 10 * time.Second
	select {
	case <-time.After(timeout):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
		return
	case closestNodeFromUserLocationResult := <-closestNodeFromUserLocationChan:
		if closestNodeFromUserLocationResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		closestNodeToResult := <-closestNodeToChan

		if err = closestNodeToResult.Err; err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
		dijkstraResultChan := make(chan db.DijkstraResult, 1)
		// -->
		go c.calculateShortestPathWorker(pipelineCtx, closestNodeFromUserLocationResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

		dijkstraResult := <-dijkstraResultChan
		if dijkstraResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(dijkstraResult.Err))
			return
		}

		nodesChan := make(chan db.Nodes, 1)
		// -->
		go c.getNodesByIdsWorker(pipelineCtx, dijkstraResult.Paths, nodesChan)

		nodesResult := <-nodesChan
		if nodesResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
		return
	case closestNodeFromResult := <-closestNodeFromChan:
		if closestNodeFromResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}

		closestNodeToResult := <-closestNodeToChan

		if err = closestNodeToResult.Err; err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
		dijkstraResultChan := make(chan db.DijkstraResult, 1)
		// -->
		go c.calculateShortestPathWorker(pipelineCtx, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

		dijkstraResult := <-dijkstraResultChan
		if dijkstraResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(dijkstraResult.Err))
			return
		}

		nodesChan := make(chan db.Nodes, 1)
		// -->
		go c.getNodesByIdsWorker(pipelineCtx, dijkstraResult.Paths, nodesChan)

		nodesResult := <-nodesChan
		if nodesResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
		return
	}
}
