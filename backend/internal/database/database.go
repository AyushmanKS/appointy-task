package database

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func InitDB() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	var err error
	DB, err = sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("Unable to open database connection: %v\n", err)
	}

	if err = DB.PingContext(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	log.Println("Successfully connected to the database.")

	createTables()
}

func createTables() {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	createURLsTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id VARCHAR(8) PRIMARY KEY,
		original_url TEXT NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		creation_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`

	createClicksTable := `
	CREATE TABLE IF NOT EXISTS clicks (
		id SERIAL PRIMARY KEY,
		url_id VARCHAR(8) NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
		clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		ip_address VARCHAR(45),
		user_agent TEXT
	);`

	ctx := context.Background()
	for _, tableSQL := range []string{createUsersTable, createURLsTable, createClicksTable} {
		if _, err := DB.ExecContext(ctx, tableSQL); err != nil {
			log.Fatalf("Unable to create table: %v\n", err)
		}
	}
	log.Println("All tables are ready.")
}
