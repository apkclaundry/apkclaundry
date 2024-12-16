package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"apkclaundry/config"
	"apkclaundry/routes"
)

func main() {
    // Initialize MongoDB connection
    if err := config.InitMongoDB(); err != nil {
        log.Fatalf("Failed to initialize MongoDB: %v", err)
    }

    // Ensure MongoDB client is disconnected properly
    defer func() {
        if err := config.Client.Disconnect(context.TODO()); err != nil {
            log.Printf("Error disconnecting MongoDB: %v", err)
        }
    }()

    // Initialize router
    router := routes.InitRoutes()

    // Define server port
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Start server
    log.Printf("Server is running on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
