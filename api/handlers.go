package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/db"
	"time"
)

const UserLocation = "My Location"

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

func (s *Server) GetPlacesAndBuildings(ctx *gin.Context) {
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get places and buildings asynchronously.
	placesChan := make(chan db.PlacesResult, 1)
	// -->
	go s.getPlacesWorker(pipelineCtx, placesChan)

	buildingsChan := make(chan db.BuildingsResult, 1)
	// -->
	go s.getBuildingsWorker(pipelineCtx, buildingsChan)

	placesResult := <-placesChan
	if err := placesResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	buildingsResult := <-buildingsChan
	if err := buildingsResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"places": placesResult.Places, "buildings": buildingsResult.Buildings})
}

func (s *Server) GetBuildings(ctx *gin.Context) {
	buildings, err := s.store.ListBuildings(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, buildings)
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

func (s *Server) GetShortestRouteByBuildingOrPlace(ctx *gin.Context) {
	var req db.RouteRequest
	err := ctx.ShouldBindJSON(&req)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Use a cancellable context to allow for early termination of the pipeline.
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	closestNodeFromChan := make(chan db.ClosestNodeResult, 1)
	closestNodeFromUserLocationChan := make(chan db.ClosestNodeToUserLocationResult, 1)

	if req.From.String == UserLocation {
		// -->
		go s.getClosestNodeByUserLocationGeom(pipelineCtx, req.FromLocation, closestNodeFromUserLocationChan)
	} else {

		geomFrom, err := s.store.GetBuildingOrPlace(pipelineCtx, req.From.String)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		// -->
		go s.getClosestNode(pipelineCtx, geomFrom.Geom, closestNodeFromChan)
	}

	// Get Building or Place for GeomTo
	geomTo, err := s.store.GetBuildingOrPlace(pipelineCtx, req.To.String)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	closestNodeToChan := make(chan db.ClosestNodeResult, 1)
	// -->
	go s.getClosestNode(pipelineCtx, geomTo.Geom, closestNodeToChan)

	timeout := 10 * time.Second
	select {
	case <-time.After(timeout):
		ctx.JSON(http.StatusRequestTimeout, gin.H{"error": "Request timed out"})
		return
	case closestNodeFromUserLocationResult := <-closestNodeFromUserLocationChan:
		if closestNodeFromUserLocationResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		closestNodeToResult := <-closestNodeToChan

		if err = closestNodeToResult.Err; err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		// Stage 3: Calculate the shortest path asynchronously.
		dijkstraResultChan := make(chan db.DijkstraResult, 1)
		// -->
		go s.calculateShortestPathWorker(pipelineCtx, closestNodeFromUserLocationResult.Node.ID, closestNodeToResult.Node.ID, dijkstraResultChan)

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
		return
	case closestNodeFromResult := <-closestNodeFromChan:
		if closestNodeFromResult.Err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		closestNodeToResult := <-closestNodeToChan

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
		return
	}
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

func (s *Server) GetShortestRouteByBuilding(ctx *gin.Context) {
	var req db.RouteRequest
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
	// -->
	go s.getClosestNode(pipelineCtx, geomFrom.BuildingCentroid, closestNodeFromChan)
	// -->
	go s.getClosestNode(pipelineCtx, geomTo.BuildingCentroid, closestNodeToChan)

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
