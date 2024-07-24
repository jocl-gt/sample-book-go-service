package main

import (
	"log"
	"readcommend/server"
)

func main() {
	app, err := server.NewServerInstance("config.local.yaml")
	if err != nil {
		log.Fatalf("Error initializing application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
