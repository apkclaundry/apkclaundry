package middleware

import (
	"apkclaundry/utils"
	"net/http"
)

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
