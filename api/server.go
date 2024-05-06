package api

import (
	"github.com/gin-gonic/gin"
	db "routing/db/sqlc"
	"routing/utils"
)

// Server all HTTP requests
type Server struct {
	store  *db.Store
	router *gin.Engine
	Graph  *utils.Graph
}

func NewServer(store *db.Store) (*Server, error) {
	server := &Server{store: store}
	server.router = gin.Default()
	err := server.router.SetTrustedProxies(nil)
	if err != nil {
		return nil, err
	}
	server.router.Static("/map", "./public")
	err = server.ConstructGraph()
	if err != nil {
		return nil, err
	}
	server.ServeRoutes()
	return server, nil
}

func (s *Server) ConstructGraph() error {
	err := s.ReadGraphIntoMemory(&gin.Context{})
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) RunServer(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
