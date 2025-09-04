// backend/cmd/server/main.go
package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	// These are your application's modules.
	"github.com/AyushmanKS/appointy-task/internal/auth"
	"github.com/AyushmanKS/appointy-task/internal/database"
	"github.com/AyushmanKS/appointy-task/internal/hub"
	"github.com/AyushmanKS/appointy-task/internal/link"
	_ "github.com/AyushmanKS/appointy-task/internal/metrics" // BLANK IMPORT to ensure metrics are registered

	// These are third-party libraries.
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

func main() {
	// --- ROBUST .env LOADING ---
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	// Go up two directories from main.go (cmd/server) to the backend folder.
	envPath := filepath.Join(basepath, "..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No backend/.env file found, relying on system environment variables.")
	}
	// --------------------------------

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("FATAL: DATABASE_URL environment variable is not set. Please ensure a backend/.env file exists with the correct content.")
	}

	database.InitDB()
	defer database.DB.Close()

	go hub.GlobalHub.Run()

	mux := http.NewServeMux()

	// --- Public Routes ---
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/r/", link.RedirectHandler)
	mux.HandleFunc("/ws", auth.WSHandler)

	// --- Protected Routes ---
	mux.Handle("/shorten", auth.JwtMiddleware(http.HandlerFunc(link.CreateLinkHandler)))
	mux.Handle("/links", auth.JwtMiddleware(http.HandlerFunc(link.GetLinksHandler)))
	mux.Handle("/analytics/", auth.JwtMiddleware(http.HandlerFunc(link.GetAnalyticsHandler)))

	// --- CORS Configuration ---
	// Using "*" for AllowedOrigins is the most flexible for local development.
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})
	handler := c.Handler(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
