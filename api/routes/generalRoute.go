package routes

func (r *Routes) GeneralRoutes() {
	r.Controller.Router.Static("/map", "./public")
	r.Controller.Router.GET("/", r.Controller.ShowMap)
	r.Controller.Router.GET("/all", r.Controller.GetPlacesAndBuildings)
	r.Controller.Router.GET("/all/route", r.Controller.GetShortestRouteByBuildingOrPlace)
	r.Controller.Router.POST("/all/route", r.Controller.GetShortestRouteByBuildingOrPlace)
	r.Controller.Router.GET("/all/search", r.Controller.FuzzyFindBuildingOrPlace)
}
