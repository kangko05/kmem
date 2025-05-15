package main

import (
	"context"
	"kmem/internal/database"
	"kmem/internal/event"
	"kmem/internal/router"
	"log"
)

const PORT = ":8000"

func main() {
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
	log.Println()

	if err := router.Setup(store).Run(PORT); err != nil {
		log.Panicf("failed to run router: %v\n", err)
	}
}
