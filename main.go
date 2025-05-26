package main

import (
	"kmem/internal/config"
	"kmem/internal/router"
	"log"
)

func main() {
	// load config
	conf, err := config.Load("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := router.Setup().Run(conf.ServerPort()); err != nil {
		log.Fatal(err)
	}
}
