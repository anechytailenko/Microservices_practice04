package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	users_api "github.com/anechytailenko/Microservices_practice04/internal/users/api"
	"github.com/anechytailenko/Microservices_practice04/internal/users/features/create_user"
	"github.com/anechytailenko/Microservices_practice04/internal/users/features/get_user"
	"github.com/anechytailenko/Microservices_practice04/internal/users/infrastructure"

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

	dbURL := os.Getenv("USERS_DB_URL")
	if dbURL == "" {
		log.Fatal("USERS_DB_URL is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	repo := infrastructure.NewPostgresRepo(db)
	createHandler := create_user.NewHandler(repo)
	getUserHandler := get_user.NewHandler(repo)

	mux := http.NewServeMux()

	shared.RegisterHealthRoutes(mux, db, "users", CommitHash)
	users_api.RegisterRoutes(mux, createHandler, getUserHandler)

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
