package routes

func (r *Routes) GeneralRoutes() {
	r.Controller.Router.GET("/", r.Controller.ServerActive)
	r.Controller.Router.GET("/all", r.Controller.GetAllEntities)
	r.Controller.Router.GET("/all/route", r.Controller.GetShortestRouteByBuildingOrPlace)
	r.Controller.Router.GET("/all/search", r.Controller.FuzzyFindBuildingOrPlaceClassroom)
}
