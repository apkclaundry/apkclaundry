package main

import (
    "context"
    "log"
    "net/http"
    "os"

    "apkclaundry/config"
    "apkclaundry/middleware"
    "apkclaundry/routes"
)

func init() {
    // Menginisialisasi MongoDB saat aplikasi dijalankan
    if err := config.InitMongoDB(); err != nil {
        log.Fatalf("Failed to initialize MongoDB: %v", err)
    }
    log.Println("MongoDB initialized successfully!")
}

func Handler(w http.ResponseWriter, r *http.Request) {
    // Inisialisasi router
    router := routes.InitRoutes()

    // Menambahkan middleware CORS
    routerWithCORS := middleware.EnableCORS(router)

    // Jalankan request melalui router
    routerWithCORS.ServeHTTP(w, r)
}

func main() {
    // Pastikan koneksi MongoDB ditutup dengan benar saat aplikasi selesai
    defer func() {
        if err := config.Client.Disconnect(context.TODO()); err != nil {
            log.Printf("Error disconnecting MongoDB: %v", err)
        }
    }()

    // Inisialisasi router
    router := routes.InitRoutes()

    // Tentukan port server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // default port
    }

    // Mulai server
    log.Printf("Server is running on port %s", port)
    if err := http.ListenAndServe(":"+port, router); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
