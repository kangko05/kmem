package main

import (
	"context"
	"kmem/internal/cache"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/queue"
	"kmem/internal/router"
	"log"

	"github.com/joho/godotenv"
)

// TODO: check dependencies - ffmpeg, ffprobe

func main() {
	// for test
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}

	// load config
	conf, err := config.Load("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	// connect postgres
	pg, err := db.Connect(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := queue.New(ctx)
	cache := cache.New(ctx)

	if err := router.Setup(pg, conf, q, cache).Run(conf.ServerPort()); err != nil {
		log.Fatal(err)
	}
}
