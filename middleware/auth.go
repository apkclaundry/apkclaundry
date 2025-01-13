package middleware

import (
	"apkclaundry/utils"
	"log"
	"net/http"
)

// EnableCORS menangani header CORS agar frontend dapat mengakses API
func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        log.Printf("Origin received: %s", origin) // Tambahkan log untuk memeriksa origin

        allowedOrigins := map[string]bool{
            "http://127.0.0.1:5500":         true,
            "http://127.0.0.1:5502":         true,
            "https://apkclaundry.github.io": true,
        }

        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Vary", "Origin")
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == http.MethodOptions {
            log.Printf("Preflight request from origin: %s", origin)
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}


// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            log.Println("Authorization header missing")
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        if len(token) < 7 || token[:7] != "Bearer " {
            log.Println("Invalid token format")
            http.Error(w, "Invalid token format", http.StatusUnauthorized)
            return
        }

        token = token[7:]
        claims, err := utils.ValidateJWT(token)
        if err != nil {
            log.Printf("Invalid token: %v", err)
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        if claims.Role != "admin" {
            log.Printf("Forbidden access for role: %s", claims.Role)
            http.Error(w, "Forbidden: Only admins can access this endpoint", http.StatusForbidden)
            return
        }

        r.Header.Set("User-ID", claims.ID)
        r.Header.Set("Username", claims.Username)
        r.Header.Set("Role", claims.Role)

        next.ServeHTTP(w, r)
    })
}



// RoleMiddleware validates user roles
func RoleMiddleware(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Header.Get("Role")
		if userRole != role {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
