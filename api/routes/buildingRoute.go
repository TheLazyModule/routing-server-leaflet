package routes

func (r *Routes) BuildingRoute() {
	r.Controller.Router.GET("/buildings/route", r.Controller.GetShortestRouteByBuilding)
	r.Controller.Router.GET("/buildings", r.Controller.GetBuildings)
}
