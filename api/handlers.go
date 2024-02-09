package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/dijkstra"
	g "routing/graph"
	"routing/utils"
)

func (s *Server) ShowMap(ctx *gin.Context) {
	//ctx.File("public/index.html")
	ctx.Redirect(http.StatusFound, "/static")
}

func (s *Server) GetPlaces(ctx *gin.Context) {
	nodes, err := s.store.ListPlaces(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodes)
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
	nodes, err := s.store.ListNodePointGeoms(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, nodes)
}

func (s *Server) GetNodePointGeomByID(ctx *gin.Context) {
	var req ReqID

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	point, err := s.store.GetNodePointGeom(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, point)
}

func (s *Server) GetNodeByID(ctx *gin.Context) {
	var req ReqID

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	nodes, err := s.store.GetNodeByID(ctx, req.ID)
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
	var req ReqID

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	edge, err := s.store.GetEdgeByID(ctx, req.ID)
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

func (s *Server) GetShortestRouteByNode(ctx *gin.Context) {
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
		return
	}

	if err = utils.ReadIntoMemory(newGraph, weights); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	paths, Distance, err := dijkstra.Dijkstra(newGraph, req.FromNodeID, req.ToNodeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	nodes, err := s.store.GetNodesByIds(ctx, paths)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": Distance, "paths": nodes})
}

func (s *Server) GetShortestRouteByPlace(ctx *gin.Context) {
	var req routeRequestByPlaceOrBuildingJSON
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	geomFrom, err := s.store.GetPlaceGeom(ctx, req.From)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	geomTo, err := s.store.GetPlaceGeom(ctx, req.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	closestNodeFrom, _ := s.store.GetClosestPointToQueryLocation(ctx, geomFrom.Location)
	closestNodeTo, _ := s.store.GetClosestPointToQueryLocation(ctx, geomTo.Location)
	fmt.Println(closestNodeTo)
	fmt.Println(closestNodeFrom)

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
		return
	}

	if err = utils.ReadIntoMemory(newGraph, weights); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	paths, Distance, err := dijkstra.Dijkstra(newGraph, closestNodeFrom.ID, closestNodeTo.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	nodes, err := s.store.GetNodesByIds(ctx, paths)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": Distance, "paths": nodes})
}

func (s *Server) GetShortestRouteByBuilding(ctx *gin.Context) {
	var req routeRequestByPlaceOrBuildingJSON
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
	closestNodeFrom, _ := s.store.GetClosestPointToQueryLocation(ctx, geomFrom.BuildingCentroidGeographic)
	closestNodeTo, _ := s.store.GetClosestPointToQueryLocation(ctx, geomTo.BuildingCentroidGeographic)
	fmt.Println(closestNodeTo)
	fmt.Println(closestNodeFrom)

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
		return
	}

	if err = utils.ReadIntoMemory(newGraph, weights); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	paths, Distance, err := dijkstra.Dijkstra(newGraph, closestNodeFrom.ID, closestNodeTo.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	nodes, err := s.store.GetNodesByIds(ctx, paths)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"distance": Distance, "paths": nodes})
}
