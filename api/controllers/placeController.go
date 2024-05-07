package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/api/utils"
	"routing/db"
)

func (c *Controller) GetPlaces(ctx *gin.Context) {
	places, err := c.store.ListPlaces(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, places)
}

func (c *Controller) GetShortestRouteByPlace(ctx *gin.Context) {
	var req db.RouteRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	geomFrom, err := c.store.GetPlaceGeom(pipelineCtx, req.From)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	geomTo, err := c.store.GetPlaceGeom(pipelineCtx, req.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	// Stage 2: Get the closest nodes asynchronously.
	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	// -->
	go c.getClosestNode(pipelineCtx, geomFrom.Geom, closestNodeFromChan)
	// -->
	go c.getClosestNode(pipelineCtx, geomTo.Geom, closestNodeToChan)

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
