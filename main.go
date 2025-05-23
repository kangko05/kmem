package main

import (
	"context"
	"fmt"
	"kmem/internal/config"
	"kmem/internal/database"
	"kmem/internal/database/query"
	"kmem/internal/event"
	"kmem/internal/router"
	"log"
	"time"
)

func main() {
	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Panicf("failed to load env: %v", err)
	// }

	// TODO: this has to be passed as env var
	const JWT_SECRET = "TODO: SECRET KEY MUST BE HANDLED WITH CARE"

	conf, err := config.Load("config.yml")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pg, err := database.Connect(ctx, &conf.Postgres)
	if err != nil {
		log.Panicf("failed to connect to postgres: %v", err)
	}
	defer pg.Close()

	log.Println("postgres connected")

	cache, err := initCache(ctx, pg)
	if err != nil {
		log.Panicf("failed to start cache: %v", err)
	}

	store := event.NewStore(ctx, pg, cache)
	go store.Run()

	log.Println("started store")
	log.Println("starting gin server...")
	fmt.Println()

	if err := router.Setup(store, conf, pg, cache, JWT_SECRET).Run(conf.GetHttpPort()); err != nil {
		log.Panicf("failed to run router: %v\n", err)
	}
}

func initCache(ctx context.Context, pg *database.Postgres) (*database.Cache, error) {
	cache := database.NewCache(ctx, time.Duration(time.Minute*10))

	usernames, err := query.QueryUsernames(pg)
	if err != nil {
		return nil, fmt.Errorf("failed to query usernames: %v", err)
	}

	for _, username := range usernames {
		userFiles, err := query.QueryUserFiles(pg, username)
		if err != nil {
			log.Printf("failed to get userfiles: %v", err)
			continue
		}

		cache.AddPermanent(username, userFiles)
	}

	return cache, nil
}
