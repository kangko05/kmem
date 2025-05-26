package main

import (
	"kmem/internal/config"
	"kmem/internal/db"
	"kmem/internal/router"
	"log"

	"github.com/joho/godotenv"
)

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

	if err := router.Setup(pg, conf).Run(conf.ServerPort()); err != nil {
		log.Fatal(err)
	}
}
