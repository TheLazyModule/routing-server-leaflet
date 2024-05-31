package controller

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"routing/config"
	db "routing/db/sqlc"
)

type Controller struct {
	store  *db.Store
	Router *gin.Engine
	Graph  *config.Graph
}

func NewController(store *db.Store, router *gin.Engine) *Controller {
	return &Controller{store: store, Router: router}
}

func (c *Controller) ConfigCORSMiddleWare() {
	c.Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://routing-web.vercel.app/"}, // Add the origin of your React app
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

}

func (c *Controller) ReadGraphIntoMemory(ctx *gin.Context) error {

	edges, err := c.store.ListEdges(ctx)
	if err != nil {
		return err
	}

	c.Graph = config.NewGraph()
	for _, edge := range edges {
		err := c.Graph.AddEdge(edge.FromNodeID, edge.ToNodeID, edge.Weight)
		if err != nil {
			return err
		}
	}
	return nil
}

func goTo() {
	// This is a dummy function
}
