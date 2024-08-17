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

func (c *Controller) GetAllEntities(ctx *gin.Context) {
	pipelineCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get places, classrooms and buildings asynchronously.
	fmt.Println("Getting places")
	placesChan := make(chan db.PlacesResult, 1)
	// -->
	go c.getPlacesWorker(pipelineCtx, placesChan)

	fmt.Println("Getting buildings")
	buildingsChan := make(chan db.BuildingsResult, 1)
	// -->
	go c.getBuildingsWorker(pipelineCtx, buildingsChan)

	classroomsChan := make(chan db.ClassroomsResult, 1)
	// -->
	fmt.Println("Getting classrooms")
	go c.getClassroomsWorker(pipelineCtx, classroomsChan)

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

	classroomsResult := <-classroomsChan
	if err := classroomsResult.Err; err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"places": placesResult.Places, "buildings": buildingsResult.Buildings, "classrooms": classroomsResult.Classrooms})
}

func (c *Controller) GetShortestRouteByBuildingOrPlace(ctx *gin.Context) {
	var req db.RouteRequest
	err := ctx.ShouldBind(&req)
	fmt.Printf("Route Request body ==> %s", req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	c.handlerBody(ctx, &req)
}

func (c *Controller) FuzzyFindBuildingOrPlaceClassroom(ctx *gin.Context) {
	var req db.SearchText
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.ErrorResponse(err))
		return
	}
	fmt.Println("Search Request ==> ", req.Text)
	result, err := c.store.FuzzyFindPlaceBuildingClassroom(ctx, req.Text)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.ErrorResponse(err))
		return
	}
	//fmt.Printf("Location Search Request ==> %v", result)

	ctx.JSON(http.StatusOK, result)
}
