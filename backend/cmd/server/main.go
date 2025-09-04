// cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AyushmanKS/appointy-task/internal/auth"
	"github.com/AyushmanKS/appointy-task/internal/database"
	"github.com/AyushmanKS/appointy-task/internal/link"

	"github.com/rs/cors" // Import the CORS library
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	// We will use only ONE router for simplicity and clarity.
	mux := http.NewServeMux()

	// --- Public Routes ---
	// These do not require authentication.
	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/r/", link.RedirectHandler)

	// --- Protected Routes ---
	// We wrap each protected handler individually with the JWT middleware.
	// This is a clearer and more standard approach than using a sub-router.
	mux.Handle("/shorten", auth.JwtMiddleware(http.HandlerFunc(link.CreateLinkHandler)))
	mux.Handle("/links", auth.JwtMiddleware(http.HandlerFunc(link.GetLinksHandler)))
	mux.Handle("/analytics/", auth.JwtMiddleware(http.HandlerFunc(link.GetAnalyticsHandler)))

	// --- CORS Configuration ---
	// This configuration allows your frontend to make authenticated requests.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://appointy-task-frontend.onrender.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Must include OPTIONS
		AllowedHeaders: []string{"Authorization", "Content-Type"},           // Must include Authorization
	})

	// Wrap our main router with the CORS handler. This ensures CORS is checked first.
	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting server on port", port)
	// Use the CORS-enabled handler to start the server.
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
