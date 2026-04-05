package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/api"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/infrastructure"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	repo := infrastructure.NewPostgresRepo(db)

	createHandler := create_meetup.NewHandler(repo)
	changeStatusHandler := change_status.NewHandler(repo)
	getMeetupHandler := get_meetup.NewHandler(repo)

	mux := api.NewRouter(createHandler, changeStatusHandler, getMeetupHandler)

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
