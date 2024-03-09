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
	graph  *utils.Graph
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	server.router = gin.Default()
	err := server.router.SetTrustedProxies(nil)
	if err != nil {
		return nil
	}
	server.router.Static("/map", "./public")
	err = server.ConstructGraph()
	if err != nil {
		return nil
	}
	server.ServeRoutes()
	return server
}

func (s *Server) ConstructGraph() error {
	newGraph, err := s.ReadGraphIntoMemory(&gin.Context{})
	if err != nil {
		return err
	}
	s.graph = newGraph
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
