package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/api/utils"
	"routing/db"
)

func (c *Controller) GetBuildings(ctx *gin.Context) {
	buildings, err := c.store.ListBuildings(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, buildings)
}

func (c *Controller) GetShortestRouteByBuilding(ctx *gin.Context) {
	var req db.RouteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Stage 1: Get building centroid geometries.
	geomFrom, err := c.store.GetBuildingCentroidGeom(pipelineCtx, req.From)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	geomTo, err := c.store.GetBuildingCentroidGeom(pipelineCtx, req.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	// Stage 2: Get the closest nodes asynchronously.
	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	// -->
	go c.getClosestNode(pipelineCtx, geomFrom.BuildingCentroid, closestNodeFromChan)
	// -->
	go c.getClosestNode(pipelineCtx, geomTo.BuildingCentroid, closestNodeToChan)

	closestNodeFromResult := <-closestNodeFromChan
	closestNodeToResult := <-closestNodeToChan
	if err = closestNodeFromResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	if err = closestNodeToResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	// Stage 3: Calculate the shortest path asynchronously.
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	// -->
	//go c.calculateShortestPathWorker(pipelineCtx, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(dijkstraResult.Err))
		return
	}

	// Stage 4: Get nodes by IDs asynchronously.
	nodesChan := make(chan db.Nodes, 1)
	// -->
	go c.getNodesByIdsWorker(pipelineCtx, dijkstraResult.Paths, nodesChan)

	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}
