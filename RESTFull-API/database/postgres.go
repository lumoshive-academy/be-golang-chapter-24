package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresDB() *sql.DB {
	connStr := "user=postgres dbname=postgres sslmode=disable password=postgres host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	return db
}
