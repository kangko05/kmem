package main

import (
	"context"
	"kmem/internal/cache"
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/queue"
	"kmem/internal/router"
	"log"
	"time"

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect postgres
	pg, err := db.Connect(ctx, conf)
	if err != nil {
		log.Fatal(err)
	}
	defer pg.Close()

	q := queue.New(ctx)

	cache := cache.New(ctx)

	q.Add(queue.CleanItems(pg, conf, cache))
	go cleanPeriod(ctx, q, pg, conf, cache)

	if err := router.Setup(pg, conf, q, cache).Run(conf.ServerPort()); err != nil {
		log.Fatal(err)
	}
}

func cleanPeriod(ctx context.Context, q *queue.Queue, pg *db.Postgres, conf *config.Config, cache *cache.Cache) {
	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			q.Add(queue.CleanItems(pg, conf, cache))
		}
	}
}
