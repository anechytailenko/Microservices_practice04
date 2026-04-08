package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	meetups_api "github.com/anechytailenko/Microservices_practice04/internal/meetups/api"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// this var will be rewritten by compiler when we would build
var CommitHash string = "unknown"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("MEETUPS_DB_URL")
	if dbURL == "" {
		log.Fatal("MEETUPS_DB_URL is not set")
	}

	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		log.Fatal("USERS_SERVICE_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := infrastructure.NewPostgresRepo(db)
	usersClient := infrastructure.NewUsersClient(usersServiceURL)

	createHandler := create_meetup.NewHandler(repo, usersClient)
	changeStatusHandler := change_status.NewHandler(repo)
	getMeetupHandler := get_meetup.NewHandler(repo)

	mux := http.NewServeMux()

	shared.RegisterHealthRoutes(mux, db, "meetups", CommitHash)
	meetups_api.RegisterRoutes(mux, createHandler, changeStatusHandler, getMeetupHandler)

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
