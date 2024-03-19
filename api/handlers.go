package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/db"
	"routing/utils"
)

func (s *Server) ShowMap(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/map")
}

func (s *Server) GetPlaces(ctx *gin.Context) {
	places, err := s.store.ListPlaces(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, places)
}

func (s *Server) GetBuildings(ctx *gin.Context) {
	nodes, err := s.store.ListBuildings(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (s *Server) GetNodes(ctx *gin.Context) {
	nodes, err := s.store.ListNodes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (s *Server) GetNodePointGeoms(ctx *gin.Context) {
	nodePointGeoms, err := s.store.ListNodePointGeoms(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodePointGeoms)
}

func (s *Server) GetNodePointGeomByID(ctx *gin.Context) {
	var req db.ReqID

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	pointGeom, err := s.store.GetNodePointGeom(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, pointGeom)
}

func (s *Server) GetNodeByID(ctx *gin.Context) {
	var req db.ReqID

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	nodeID, err := s.store.GetNodeByID(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodeID)
}

func (s *Server) GetEdges(ctx *gin.Context) {
	edges, err := s.store.ListEdges(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, edges)
}

func (s *Server) GetShortestRouteByNodes(ctx *gin.Context) {
	var req db.RouteRequestByID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	nodesChan := make(chan db.Nodes, 1)

	go s.calculateShortestPath(ctx, req.FromNodeID, req.ToNodeID, dijkstraResultChan)

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	go s.getNodesByIds(ctx, dijkstraResult.Paths, nodesChan)
	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(nodesResult.Err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}

func (s *Server) GetShortestRouteByPlace(ctx *gin.Context) {
	var req db.RouteRequestByBuildingJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Stage 1: Get building centroid geometries.
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
	go s.getClosestNode(pipelineCtx, geomFrom.Geom, closestNodeFromChan)
	go s.getClosestNode(pipelineCtx, geomTo.Geom, closestNodeToChan)

	closestNodeFromResult := <-closestNodeFromChan
	closestNodeToResult := <-closestNodeToChan
	if closestNodeFromResult.Err != nil || closestNodeToResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Stage 3: Calculate the shortest path asynchronously.
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	go s.calculateShortestPath(pipelineCtx, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	// Stage 4: Get nodes by IDs asynchronously.
	nodesChan := make(chan db.Nodes, 1)
	go s.getNodesByIds(pipelineCtx, dijkstraResult.Paths, nodesChan)

	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}

func (s *Server) GetShortestRouteByBuilding(ctx *gin.Context) {
	var req db.RouteRequestByBuildingJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Stage 1: Get building centroid geometries.
	geomFrom, err := s.store.GetBuildingCentroidGeom(pipelineCtx, req.From)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	geomTo, err := s.store.GetBuildingCentroidGeom(pipelineCtx, req.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Stage 2: Get the closest nodes asynchronously.
	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	go s.getClosestNode(pipelineCtx, geomFrom.BuildingCentroid, closestNodeFromChan)
	go s.getClosestNode(pipelineCtx, geomTo.BuildingCentroid, closestNodeToChan)

	closestNodeFromResult := <-closestNodeFromChan
	closestNodeToResult := <-closestNodeToChan
	if closestNodeFromResult.Err != nil || closestNodeToResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Stage 3: Calculate the shortest path asynchronously.
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	go s.calculateShortestPath(pipelineCtx, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	// Stage 4: Get nodes by IDs asynchronously.
	nodesChan := make(chan db.Nodes, 1)
	go s.getNodesByIds(pipelineCtx, dijkstraResult.Paths, nodesChan)

	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}

func (s *Server) getClosestNode(ctx context.Context, centroid interface{}, resultChan chan<- db.ClosestNodeResult) {
	defer close(resultChan)
	node, err := s.store.GetClosestPointToQueryLocation(ctx, centroid)
	select {
	case resultChan <- db.ClosestNodeResult{Node: node, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) calculateShortestPath(ctx context.Context, fromID, toID int64, resultChan chan<- db.DijkstraResult) {
	defer close(resultChan)
	paths, distance, err := utils.Dijkstra(s.Graph, fromID, toID)
	select {
	case resultChan <- db.DijkstraResult{Paths: paths, Distance: distance, Err: err}:
	case <-ctx.Done():
	}
}

func (s *Server) getNodesByIds(ctx context.Context, ids []int64, resultChan chan<- db.Nodes) {
	defer close(resultChan)
	nodes, err := s.store.GetNodesByIds(ctx, ids)
	select {
	case resultChan <- db.Nodes{Nodes: nodes, Err: err}:
	case <-ctx.Done():
	}
}
