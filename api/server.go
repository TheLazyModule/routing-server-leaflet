package api

import (
	"github.com/gin-gonic/gin"
	controller "routing/api/controllers"
	"routing/api/routes"
	db "routing/db/sqlc"
)

// Server all HTTP requests
type Server struct {
	Controller *controller.Controller
	Routes     *routes.Routes
}

func NewServer(store *db.Store) (*Server, error) {
	_controller := controller.NewController(store, gin.Default())
	_routes := routes.NewRoutes(_controller)
	server := &Server{_controller, _routes}
	server.Controller.Router.Static("/map", "./public")
	err := server.ConstructGraph()
	if err != nil {
		return nil, err
	}
	server.ServeRoutes()
	return server, nil
}

func (s *Server) ServeRoutes() {
	s.Routes.BuildingRoute()
	s.Routes.NodeRoute()
	s.Routes.EdgeRoute()
	s.Routes.ClassroomRoutes()
	s.Routes.GeneralRoutes()
}

func (s *Server) ConstructGraph() error {
	err := s.Controller.ReadGraphIntoMemory(&gin.Context{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) RunServer(address string) error {
	return s.Controller.Router.Run(address)
}
