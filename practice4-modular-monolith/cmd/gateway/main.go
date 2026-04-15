package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/gateway/reverseproxy"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/web"
)

func main() {
	var CommitHash string
	if envHash := os.Getenv("COMMIT_HASH"); envHash != "" {
		CommitHash = envHash
	}

	meetupsURL := os.Getenv("MEETUPS_SERVICE_URL")
	if meetupsURL == "" {
		log.Fatal("MEETUPS_SERVICE_URL is not set")
	}

	usersURL := os.Getenv("USERS_SERVICE_URL")
	if usersURL == "" {
		log.Fatal("USERS_SERVICE_URL is not set")
	}

	notificationsURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notificationsURL == "" {
		log.Fatal("NOTIFICATION_SERVICE_URL is not set")
	}

	workflowURL := os.Getenv("WORKFLOW_SERVICE_URL")
	if workflowURL == "" {
		log.Fatal("WORKFLOW_SERVICE_URL is not set")
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

	notificationsProxy, err := reverseproxy.NewProxy(notificationsURL)
	if err != nil {
		log.Fatalf("Failed to create notifications proxy: %v", err)
	}

	workflowProxy, err := reverseproxy.NewProxy(workflowURL)
	if err != nil {
		log.Fatalf("Failed to create workflow proxy: %v", err)
	}

	mux := http.NewServeMux()

	web.RegisterHealthRoutes(mux, nil, "gateway", CommitHash)

	mux.Handle("/users", http.StripPrefix("/users", usersProxy))
	mux.Handle("/users/", http.StripPrefix("/users", usersProxy))

	mux.Handle("/meetups", http.StripPrefix("/meetups", meetupsProxy))
	mux.Handle("/meetups/", http.StripPrefix("/meetups", meetupsProxy))

	mux.Handle("/notifications/", http.StripPrefix("/notifications", notificationsProxy))

	mux.Handle("/workflows/", http.StripPrefix("/workflows", workflowProxy))

	finalHandler := reverseproxy.CorrelationID(mux)

	log.Printf("API Gateway is starting on port %s...", port)

	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		log.Fatalf("Gateway server failed: %v", err)
	}
}
