package main

import (
	"context"
	"fmt"
	"kmem/internal/database"
	"kmem/internal/event"
	"kmem/internal/router"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Panicf("failed to load env: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pg, err := database.Connect(ctx)
	if err != nil {
		log.Panicf("failed to connect to postgres: %v", err)
	}
	defer pg.Close()

	log.Println("postgres connected")

	store := event.NewStore(ctx, pg)
	go store.Run()

	log.Println("started store")
	log.Println("starting gin server...")
	fmt.Println()

	if err := router.Setup(store, pg).Run(os.Getenv("PORT")); err != nil {
		log.Panicf("failed to run router: %v\n", err)
	}
}
