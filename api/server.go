package api

import (
	"github.com/gin-contrib/cors"
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
	// Initialize the controller
	_controller := controller.NewController(store, gin.Default())
	// Initialize the routes
	_routes := routes.NewRoutes(_controller)
	// Initialize the server
	server := &Server{_controller, _routes}
	// Configure CORS settings
	server.Controller.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://routing-web.vercel.app/"}, // Add the origin of your React app
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// Construct the graph
	err := server.ConstructGraph()
	if err != nil {
		return nil, err
	}
	server.SetRoutes()
	return server, nil
}

func (s *Server) SetRoutes() {
	s.Routes.GeneralRoutes()
	s.Routes.BuildingRoute()
	s.Routes.PlaceRoute()
	s.Routes.NodeRoute()
	s.Routes.EdgeRoute()
	//s.Routes.ClassroomRoutes()
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
