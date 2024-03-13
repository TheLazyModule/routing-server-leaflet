package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"routing/api"
	db "routing/db/sqlc"
	"routing/utils"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
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
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Cannot initialize Server")
	}

	//paths, distance, err := utils.Dijkstra(server.Graph, 150, 700)
	//if err != nil {
	//	return
	//}
	//fmt.Println(paths, distance)

	err = server.RunServer(config.ServerAddress)
	if err != nil {
		log.Fatal("Cannot RunServer Server")
	}

}
