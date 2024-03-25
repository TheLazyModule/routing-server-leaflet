package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/db"
)

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

	// -->
	go s.calculateShortestPathWorker(ctx, req.FromNodeID, req.ToNodeID, dijkstraResultChan)

	dijkstraResult := <-dijkstraResultChan
	if dijkstraResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(dijkstraResult.Err))
		return
	}

	// -->
	go s.getNodesByIdsWorker(ctx, dijkstraResult.Paths, nodesChan)
	nodesResult := <-nodesChan
	if nodesResult.Err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(nodesResult.Err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": dijkstraResult.Distance, "paths": nodesResult.Nodes})
}
