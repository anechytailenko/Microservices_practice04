package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
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

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "UP",
		})
	})

	// POST  /meetups
	mux.HandleFunc("POST /meetups", func(w http.ResponseWriter, r *http.Request) {
		var cmd create_meetup.Command

		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		id, err := createHandler.Handle(r.Context(), cmd)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"id": string(id)})
	})

	// PATCH  /meetups/{id}/status)
	mux.HandleFunc("PATCH /meetups/{id}/status", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var reqBody struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		cmd := change_status.Command{
			MeetupID: id,
			Status:   reqBody.Status,
		}

		if err := changeStatusHandler.Handle(r.Context(), cmd); err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteNoContent(w)
	})

	// GET /meetups/{id}
	mux.HandleFunc("GET /meetups/{id}", func(w http.ResponseWriter, r *http.Request) {
		q := get_meetup.Query{
			MeetupID: r.PathValue("id"),
		}

		dto, err := getMeetupHandler.Handle(r.Context(), q)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusOK, dto)
	})

	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
