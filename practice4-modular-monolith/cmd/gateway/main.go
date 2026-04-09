package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/gateway/reverseproxy"
)

// this var will be rewritten by compiler when we would build
var CommitHash string = "unknown"

func main() {
	meetupsURL := os.Getenv("MEETUPS_SERVICE_URL")
	if meetupsURL == "" {
		log.Fatal("MEETUPS_SERVICE_URL is not set")
	}

	usersURL := os.Getenv("USERS_SERVICE_URL")
	if usersURL == "" {
		log.Fatal("USERS_SERVICE_URL is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	meetupsProxy, err := reverseproxy.NewProxy(meetupsURL)
	if err != nil {
		log.Fatalf("Failed to create meetups proxy: %v", err)
	}

	usersProxy, err := reverseproxy.NewProxy(usersURL)
	if err != nil {
		log.Fatalf("Failed to create users proxy: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "UP",
			"service": "gateway",
			"commit":  CommitHash,
		})
	})

	mux.Handle("/users/", http.StripPrefix("/users", usersProxy))
	mux.Handle("POST /users", http.StripPrefix("/users", usersProxy))

	mux.Handle("/meetups/", http.StripPrefix("/meetups", meetupsProxy))
	mux.Handle("POST /meetups", http.StripPrefix("/meetups", meetupsProxy))
	mux.Handle("PATCH /meetups/", http.StripPrefix("/meetups", meetupsProxy))

	finalHandler := reverseproxy.CorrelationID(mux)

	log.Printf("API Gateway is starting on port %s...", port)

	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatalf("Gateway server failed: %v", err)
	}
}
