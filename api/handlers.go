package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/db"
	"routing/utils"
	"sync"
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

	var wg sync.WaitGroup
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	nodesChan := make(chan db.Nodes, 1)

	go func() {
		defer wg.Done()
		paths, distance, err := utils.Dijkstra(s.Graph, req.FromNodeID, req.ToNodeID)
		dijkstraResultChan <- db.DijkstraResult{Paths: paths, Distance: distance, Err: err}
	}()

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	go func() {
		defer wg.Done()
		nodes, err := s.store.GetNodesByIds(ctx, dijkstraResult.Paths)
		nodesChan <- db.Nodes{Nodes: nodes, Err: err}
	}()
	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(nodesResult.Err))
		return
	}

	wg.Wait()
	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}

func (s *Server) GetShortestRouteByPlace(ctx *gin.Context) {
	var req db.RouteRequestByPlaceOrBuildingJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	geomFrom, err := s.store.GetPlaceGeom(ctx, req.From.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	geomTo, err := s.store.GetPlaceGeom(ctx, req.To.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	nodesChan := make(chan db.Nodes, 1)

	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		closestNodeFrom, err := s.store.GetClosestPointToQueryLocation(ctx, geomFrom.Geom)
		closestNodeFromChan <- db.ClosestNodeResult{Node: closestNodeFrom, Err: err}
	}()
	go func() {
		defer wg.Done()
		closestNodeTo, err := s.store.GetClosestPointToQueryLocation(ctx, geomTo.Geom)
		closestNodeToChan <- db.ClosestNodeResult{Node: closestNodeTo, Err: err}
	}()

	closestNodeFromResult := <-closestNodeFromChan
	closestNodeToResult := <-closestNodeToChan

	if closestNodeFromResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if closestNodeToResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go func() {
		defer wg.Done()
		paths, distance, err := utils.Dijkstra(s.Graph, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID)
		dijkstraResultChan <- db.DijkstraResult{Paths: paths, Distance: distance, Err: err}
	}()

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	go func() {
		defer wg.Done()
		nodes, err := s.store.GetNodesByIds(ctx, dijkstraResult.Paths)
		nodesChan <- db.Nodes{Nodes: nodes, Err: err}
	}()
	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	wg.Wait()
	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}

func (s *Server) GetShortestRouteByBuilding(ctx *gin.Context) {
	var req db.RouteRequestByPlaceOrBuildingJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	geomFrom, err := s.store.GetBuildingCentroidGeom(ctx, req.From)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	geomTo, err := s.store.GetBuildingCentroidGeom(ctx, req.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	dijkstraResultChan := make(chan db.DijkstraResult, 1)
	nodesChan := make(chan db.Nodes, 1)

	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		closestNodeFrom, err := s.store.GetClosestPointToQueryLocation(ctx, geomFrom.BuildingCentroid)
		closestNodeFromChan <- db.ClosestNodeResult{Node: closestNodeFrom, Err: err}
	}()
	go func() {
		defer wg.Done()
		closestNodeTo, err := s.store.GetClosestPointToQueryLocation(ctx, geomTo.BuildingCentroid)
		closestNodeToChan <- db.ClosestNodeResult{Node: closestNodeTo, Err: err}
	}()

	closestNodeFromResult := <-closestNodeFromChan
	closestNodeToResult := <-closestNodeToChan

	if closestNodeFromResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if closestNodeToResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	go func() {
		defer wg.Done()
		paths, distance, err := utils.Dijkstra(s.Graph, closestNodeFromResult.Node.ID, closestNodeToResult.Node.ID)
		dijkstraResultChan <- db.DijkstraResult{Paths: paths, Distance: distance, Err: err}
	}()

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	go func() {
		defer wg.Done()
		nodes, err := s.store.GetNodesByIds(ctx, dijkstraResult.Paths)
		nodesChan <- db.Nodes{Nodes: nodes, Err: err}
	}()
	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	wg.Wait()
	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}
