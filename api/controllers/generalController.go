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

func (c *Controller) ShowMap(ctx *gin.Context) {
	ctx.Redirect(http.StatusFound, "/map")
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
	var reqJSON db.RouteRequestJSON
	contentType := ctx.GetHeader("Content-Type")
	// Query Params
	if contentType == "application/json" {
		err := ctx.ShouldBind(&req)
		fmt.Println(req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}
		c.handlerBody(ctx, &req)
	} else {
		// Static web request
		err := ctx.ShouldBindJSON(&reqJSON)
		fmt.Println(reqJSON)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
			return
		}
		c.handlerBody(ctx, &reqJSON)
	}
}
