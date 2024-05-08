package routes

func (r *Routes) GeneralRoutes() {
	r.Controller.Router.GET("/", r.Controller.ServerActive)
	r.Controller.Router.GET("/all", r.Controller.GetPlacesAndBuildings)
	r.Controller.Router.GET("/all/route", r.Controller.GetShortestRouteByBuildingOrPlace)
	r.Controller.Router.GET("/all/search", r.Controller.FuzzyFindBuildingOrPlace)
	r.Controller.Router.GET("all/location/search", r.Controller.GetLocationByName)
}
