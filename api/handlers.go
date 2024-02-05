package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetNodes(ctx *gin.Context) {
	nodes, err := s.store.ListNodes(ctx)
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

	//_, err := s.store.GetNodeAndEdges(ctx, req.ToNodeID)
	//_, err := s.store.GetNodeAndEdges(ctx, req.FromNodeID)
	//nodes, err := s.store.ListNodes(ctx)
	edges, err := s.store.ListEdges(ctx)
	//weights, err := s.store.ListWeights(ctx)
	for _, e := range edges {
		fmt.Println(e)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
