package api

func (s *Server) ServeRoutes() {
	s.router.GET("/", s.Map)
	s.router.GET("/buildings", s.GetBuildings)
	s.router.GET("/places", s.GetPlaces)
	s.router.GET("/nodes", s.GetNodes)
	s.router.GET("/nodes/:id", s.GetNodeByID)
	s.router.GET("/nodes/geom/:id", s.GetNodePointGeomByID)
	s.router.GET("/nodes/geoms", s.GetNodePointGeoms)
	s.router.GET("/edges", s.GetEdges)
	s.router.GET("/edges/:id", s.GetEdgeByID)
	s.router.GET("/weights", s.GetWeights)
	s.router.POST("/route", s.GetShortestRoute)
}
