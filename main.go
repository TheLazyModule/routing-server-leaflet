package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"routing/api"
	"routing/config"
	db "routing/db/sqlc"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	configEnv, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot Load configurations")
	}

	conn, err := pgxpool.New(context.Background(), configEnv.DBUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	} else {
		fmt.Println("Database Connected!")
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Cannot initialize Server")
	}

	err = server.RunServer(configEnv.ServerAddress)
	if err != nil {
		log.Fatal("Cannot RunServer Server")
	}

}
