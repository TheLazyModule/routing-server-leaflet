package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/dijkstra"
	g "routing/graph"
	"routing/utils"
)

func (s *Server) GetNodes(ctx *gin.Context) {
	nodes, err := s.store.ListNodes(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (s *Server) GetNodeByID(ctx *gin.Context) {
	var req int64

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	nodes, err := s.store.GetNodeByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (s *Server) GetEdges(ctx *gin.Context) {
	edges, err := s.store.ListEdges(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, edges)

}

func (s *Server) GetEdgeByID(ctx *gin.Context) {
	var req int64

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	edge, err := s.store.GetEdgeByID(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, edge)
}

func (s *Server) GetWeights(ctx *gin.Context) {
	weights, err := s.store.ListWeights(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, weights)
}

func (s *Server) GetShortestRoute(ctx *gin.Context) {
	var req routeRequestByID
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	edges, err := s.store.ListEdges(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	weights, err := s.store.ListWeights(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	newGraph := g.NewGraph()
	if err = utils.ReadIntoMemory(newGraph, edges); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	if err = utils.ReadIntoMemory(newGraph, weights); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	paths, err := dijkstra.Dijkstra(newGraph, req.FromNodeID, req.ToNodeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	nodes, err := s.store.GetNodesByIds(ctx, paths)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, gin.H{"paths": nodes})
}
