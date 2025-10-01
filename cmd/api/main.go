package main

import (
	"context"
	"log"

	"katseye/internal/infrastructure/config"
)

func main() {
	ctx := context.Background()

	log.Println("startup: initializing application")

	app, err := config.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		log.Println("shutdown: releasing resources")
		if closeErr := app.Close(ctx); closeErr != nil {
			log.Printf("error closing application resources: %v", closeErr)
		}
	}()

	log.Println("startup: boot sequence complete, starting HTTP server")

	if err := app.RunHTTPServer(); err != nil {
		log.Fatal(err)
	}
}
