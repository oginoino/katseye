package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"katseye/internal/infrastructure/config"
)

func main() {
	app, err := config.Initialize()
	if err != nil {
		log.Fatalf("initialization failed: %v", err)
	}

	go func() {
		if err := app.RunHTTPServer(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.HTTPServer().Shutdown(ctx); err != nil {
		log.Printf("server shutdown failed: %v", err)
	}

	if err := app.Close(ctx); err != nil {
		log.Printf("closing resources failed: %v", err)
	}

	log.Println("shutdown complete")
}
