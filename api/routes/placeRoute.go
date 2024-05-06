package routes

func (r *Routes) PlaceRoute() {
	r.Controller.Router.POST("/places/route", r.Controller.GetShortestRouteByPlace)
	r.Controller.Router.GET("/places", r.Controller.GetPlaces)
}
