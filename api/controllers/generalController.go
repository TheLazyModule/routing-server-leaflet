package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"routing/api/utils"
	"routing/db"
)

var UserLocation string = "My Location"

func (c *Controller) ServerActive(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Server is up and Running"})
}

func (c *Controller) GetPlacesAndBuildings(ctx *gin.Context) {
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get places and buildings asynchronously.
	placesChan := make(chan db.PlacesResult, 1)
	// -->
	go c.getPlacesWorker(pipelineCtx, placesChan)

	buildingsChan := make(chan db.BuildingsResult, 1)
	// -->
	go c.getBuildingsWorker(pipelineCtx, buildingsChan)

	placesResult := <-placesChan
	if err := placesResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	buildingsResult := <-buildingsChan
	if err := buildingsResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"places": placesResult.Places, "buildings": buildingsResult.Buildings})
}

func (c *Controller) GetShortestRouteByBuildingOrPlace(ctx *gin.Context) {
	var req db.RouteRequest
	err := ctx.ShouldBind(&req)
	fmt.Println(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	c.handlerBody(ctx, &req)
}

func (c *Controller) FuzzyFindBuildingOrPlace(ctx *gin.Context) {
	var req db.SearchText
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
	}
	fmt.Println(req.Text)

	result, err := c.store.FuzzyFindPlaceOrBuilding(ctx, req.Text)
	if err != nil {
		return
	}
	ctx.JSON(http.StatusOK, result)
}
