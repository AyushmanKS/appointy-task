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

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/r/", link.RedirectHandler)

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/shorten", link.CreateLinkHandler)
	protectedMux.HandleFunc("/links", link.GetLinksHandler)
	protectedMux.HandleFunc("/analytics/", link.GetAnalyticsHandler)

	// Apply JWT middleware to protected routes
	mux.Handle("/shorten", auth.JwtMiddleware(protectedMux))
	mux.Handle("/links", auth.JwtMiddleware(protectedMux))
	mux.Handle("/analytics/", auth.JwtMiddleware(protectedMux))

	// --- NEW CORS CONFIGURATION ---
	// This creates a CORS handler that allows requests from your frontend.
	c := cors.New(cors.Options{
		// IMPORTANT: For better security, we specify the exact origin.
		AllowedOrigins: []string{"https://appointy-task-frontend.onrender.com"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	// Wrap your main router with the CORS handler
	handler := c.Handler(mux)
	// -------------------------------

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting server on port", port)
	// Use the new 'handler' which has CORS enabled
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
