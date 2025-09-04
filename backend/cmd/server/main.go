package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AyushmanKS/appointy-task/internal/auth"
	"github.com/AyushmanKS/appointy-task/internal/database"
	"github.com/AyushmanKS/appointy-task/internal/link"
)

func main() {
	database.InitDB()
	defer database.DB.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/register", auth.RegisterHandler)
	mux.HandleFunc("/login", auth.LoginHandler)
	mux.HandleFunc("/r/", link.RedirectHandler)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/shorten", link.CreateLinkHandler)
	protectedMux.HandleFunc("/analytics/", link.GetAnalyticsHandler)

	mux.Handle("/shorten", auth.JwtMiddleware(protectedMux))
	mux.Handle("/analytics/", auth.JwtMiddleware(protectedMux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
