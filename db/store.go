package db

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	server.Cors()
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

func (s *Server) Cors() {
	// CORS middleware
	s.router.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Set your allowed origin
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusOK)
			return
		}
		ctx.Next()
	})

}
