// backend/internal/auth/ws_handler.go
package auth

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/AyushmanKS/appointy-task/internal/hub"
	"github.com/golang-jwt/jwt/v5" // <-- THIS LINE IS NOW CORRECT
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local development
	},
}

// WSHandler handles WebSocket connection requests.
func WSHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	client := &hub.Client{UserID: userID, Conn: conn}

	hub.GlobalHub.RegisterClient(client)

	// This goroutine keeps the connection alive and handles unregistering.
	go func() {
		defer func() {
			hub.GlobalHub.UnregisterClient(client)
		}()
		for {
			// Read messages to detect when the client closes the connection.
			if _, _, err := conn.ReadMessage(); err != nil {
				break // Exit loop on error, which triggers the defer.
			}
		}
	}()
}
