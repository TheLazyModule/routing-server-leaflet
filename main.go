package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"os/signal"
	"routing/api"
	"routing/config"
	db "routing/db/sqlc"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	fmt.Println("Loading configurations...")
	configEnv, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configurations:", err)
	}
	fmt.Println("Configurations loaded")

	fmt.Println("Connecting to database...")
	conn, err := pgxpool.New(context.Background(), configEnv.DBUrl)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	} else {
		fmt.Println("Database Connected!")
	}

	acquire, err := conn.Acquire(context.Background())
	fmt.Println("Acquiring connection to deallocate prepared statements...")
	if err != nil {
		log.Fatalf("Unable to acquire connection: %v\n", err)
	}
	defer acquire.Release()

	fmt.Println("Deallocating all prepared statements...")
	err = acquire.Conn().DeallocateAll(context.Background())
	if err != nil {
		log.Fatalf("Unable to deallocate prepared statements: %v\n", err)
	}
	fmt.Println("Deallocated all prepared statements!")

	fmt.Println("Creating new store...")
	store := db.NewStore(conn)

	fmt.Println("Initializing server...")
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("Cannot initialize server:", err)
	}
	fmt.Println("Server initialized")

	// Graceful shutdown
	fmt.Println("Setting up graceful shutdown...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Shutting down gracefully...")
		conn.Close()
		os.Exit(0)
	}()

	fmt.Println("Running Server on address:", configEnv.ServerAddress)
	err = server.RunServer(configEnv.ServerAddress)

	if err != nil {
		log.Fatal("Cannot run server:", err)
	}

	fmt.Println("Server is running on address:", configEnv.ServerAddress)
}
