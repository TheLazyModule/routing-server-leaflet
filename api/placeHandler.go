package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/db"
)

func (s *Server) GetPlaces(ctx *gin.Context) {
	places, err := s.store.ListPlaces(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, places)
}

func (s *Server) GetShortestRouteByPlace(ctx *gin.Context) {
	var req db.RouteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	geomFrom, err := s.store.GetPlaceGeom(pipelineCtx, req.From.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	geomTo, err := s.store.GetPlaceGeom(pipelineCtx, req.To.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Stage 2: Get the closest nodes asynchronously.
	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	// -->
	go s.getClosestNode(pipelineCtx, geomFrom.Geom, closestNodeFromChan)
	// -->
	go s.getClosestNode(pipelineCtx, geomTo.Geom, closestNodeToChan)

	closestNodeFromResult := <-closestNodeFromChan
	closestNodeToResult := <-closestNodeToChan
	if err = closestNodeFromResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err = closestNodeToResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Stage 3: Calculate the shortest path asynchronously.
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	// -->
	go s.calculateShortestPathWorker(pipelineCtx, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	// Stage 4: Get nodes by IDs asynchronously.
	nodesChan := make(chan db.Nodes, 1)
	// -->
	go s.getNodesByIdsWorker(pipelineCtx, dijkstraResult.Paths, nodesChan)

	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}
