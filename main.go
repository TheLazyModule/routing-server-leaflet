package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"routing/api"
	db "routing/db/sqlc"
	"routing/dijkstra"
	rg "routing/graph"
	"routing/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load configurations")
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Database Connected!")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.RunServer(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot RunServer Server")
	}

}

func Dk() {

	var newGraph = rg.NewGraph()

	type edges struct {
		fromNode, toNode string
		weight           float64
	}

	var edgesContainer = []edges{
		{"X", "A", 7},
		{"A", "B", 3},
		{"A", "D", 4},
		{"X", "B", 2},
		{"X", "C", 3},
		{"B", "H", 5},
		{"B", "D", 4},
		{"C", "L", 2},
		{"D", "F", 1},
		{"X", "E", 4},
		{"F", "H", 3},
		{"G", "Y", 2},
		{"G", "H", 2},
		{"I", "J", 6},
		{"I", "L", 4},
		{"I", "K", 4},
		{"K", "Y", 5},
		{"J", "L", 1},
	}

	for _, edge := range edgesContainer {
		newGraph.AddEdge(edge.fromNode, edge.toNode, edge.weight)
	}

	//fmt.Println(newGraph.GetWeights())
	//fmt.Println(newGraph.GetEdges())
	path, err := dijkstra.Dijkstra(newGraph, "X", "K")
	if err == nil {
		fmt.Println(path)
	}
}
