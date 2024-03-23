package api

func (s *Server) ServeRoutes() {
	s.router.GET("/", s.ShowMap)
	s.router.POST("/places/route", s.GetShortestRouteByPlace)
	s.router.POST("/buildings/route", s.GetShortestRouteByBuilding)
	s.router.GET("/buildings", s.GetBuildings)
	s.router.GET("/places", s.GetPlaces)
	s.router.GET("/all", s.GetPlacesAndBuildings)
	s.router.POST("/all/route", s.GetShortestRouteByBuildingOrPlace)
	s.router.GET("/nodes", s.GetNodes)
	s.router.GET("/nodes/:id", s.GetNodeByID)
	s.router.GET("/nodes/geoms/:id", s.GetNodePointGeomByID)
	s.router.GET("/nodes/geoms", s.GetNodePointGeoms)
	s.router.GET("/edges", s.GetEdges)
	s.router.POST("/nodes/route", s.GetShortestRouteByNodes)
}
