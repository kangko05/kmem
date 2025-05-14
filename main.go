package main

import (
	"kmem/internal/router"
	"log"
)

const PORT = ":8000"

func main() {
	if err := router.Setup().Run(PORT); err != nil {
		log.Panicf("failed to run router: %v\n", err)
	}
}
