package api

import (
	"github.com/gin-gonic/gin"
	db "routing/db/sqlc"
)

// Server all HTTP requests
type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	server.router = gin.Default()
	err := server.router.SetTrustedProxies(nil)
	server.ServeRoutes()
	if err != nil {
		return nil
	}
	return server
}

func (s *Server) RunServer(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
