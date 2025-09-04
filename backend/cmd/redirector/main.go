package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AyushmanKS/appointy-task/internal/database"
	"github.com/AyushmanKS/appointy-task/internal/link"
	_ "github.com/AyushmanKS/appointy-task/internal/metrics"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	envPath := filepath.Join(basepath, "..", "..", ".env")
	godotenv.Load(envPath)

	if os.Getenv("DATABASE_URL") == "" {
		log.Fatal("FATAL: DATABASE_URL is not set.")
	}

	database.InitDB()
	defer database.DB.Close()

	r := chi.NewRouter()
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/r/{id}", link.RedirectHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting Redirector Microservice on port", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
