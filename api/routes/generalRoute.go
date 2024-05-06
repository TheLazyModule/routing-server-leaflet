package routes

func (r *Routes) GeneralRoutes() {
	r.Controller.Router.Static("/public", "./public")
	r.Controller.Router.POST("/map", r.Controller.ShowMap)
	r.Controller.Router.POST("/places/route", r.Controller.GetShortestRouteByPlace)
	r.Controller.Router.POST("all", r.Controller.GetPlacesAndBuildings)
	r.Controller.Router.GET("/", r.Controller.ShowMap)
	r.Controller.Router.GET("/all/route", r.Controller.GetShortestRouteByBuildingOrPlace)
}
