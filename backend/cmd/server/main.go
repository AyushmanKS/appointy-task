package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AyushmanKS/appointy-task/internal/auth"
	"github.com/AyushmanKS/appointy-task/internal/database"
	"github.com/AyushmanKS/appointy-task/internal/hub"
	"github.com/AyushmanKS/appointy-task/internal/link"
	_ "github.com/AyushmanKS/appointy-task/internal/metrics"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	envPath := filepath.Join(basepath, "..", "..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Println("No backend/.env file found.")
	}

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("FATAL: DATABASE_URL is not set.")
	}

	database.InitDB()
	defer database.DB.Close()

	go hub.GlobalHub.Run()

	r := chi.NewRouter()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
	r.Use(corsHandler.Handler)

	r.Group(func(r chi.Router) {
		r.Get("/metrics", promhttp.Handler().ServeHTTP)
		r.Post("/register", auth.RegisterHandler)
		r.Post("/login", auth.LoginHandler)
		r.Get("/r/{id}", link.RedirectHandler)
		r.Get("/ws", auth.WSHandler)
	})

	r.Group(func(r chi.Router) {
		r.Use(auth.JwtMiddleware)
		r.Post("/shorten", link.CreateLinkHandler)
		r.Get("/links", link.GetLinksHandler)
		r.Get("/analytics/{id}", link.GetAnalyticsHandler)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
