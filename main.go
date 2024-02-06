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

}
