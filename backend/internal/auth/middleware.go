// internal/auth/middleware.go
package auth

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Use a custom type for the context key to avoid collisions.
type contextKey string

// UserContextKey is an exported key for accessing the user ID in the context.
const UserContextKey = contextKey("userID")

// JwtMiddleware protects routes by validating the JWT.
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &jwt.RegisteredClaims{}
		jwtKey := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			http.Error(w, "Invalid user ID in token", http.StatusInternalServerError)
			return
		}

		// Add user ID to the request context using our exported key.
		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- NEW HELPER FUNCTION ---
// GetUserIDFromContext safely retrieves the user ID from the context.
// It returns the user ID and true if successful, or 0 and false if not.
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserContextKey).(int)
	return userID, ok
}
