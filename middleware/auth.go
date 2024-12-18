package middleware

import (
	"apkclaundry/utils"
	"net/http"
)

// EnableCORS menangani header CORS agar frontend dapat mengakses API
func EnableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // Daftar origin yang diperbolehkan
        allowedOrigins := map[string]bool{
            "http://127.0.0.1:5500":                      true,
            "https://proyek3-pos.github.io/laundrypos-fe": true,
            "https://proyek3-pos.github.io/swagger":       true,
            "https://apkclaundry.github.io/":          true,
        }

        // Periksa apakah origin dalam daftar yang diizinkan
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }

        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Tangani preflight request (OPTIONS)
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Lanjutkan ke handler berikutnya
        next.ServeHTTP(w, r)
    })
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Ensure the token starts with "Bearer "
		if len(token) < 7 || token[:7] != "Bearer " {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Extract the token after "Bearer "
		token = token[7:]

		// Validate the token
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user information to the request context for further use
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
