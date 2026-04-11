package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/notifications/api"
	"github.com/anechytailenko/Microservices_practice04/internal/notifications/features/get_notifications"
	"github.com/anechytailenko/Microservices_practice04/internal/notifications/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var CommitHash string = "unknown"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("NOTIFICATIONS_DB_URL")
	if dbURL == "" {
		log.Fatal("NOTIFICATIONS_DB_URL is not set")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	subscriber, msgsChan, err := rabbitmq.NewSubscriber(
		rabbitMQURL,
		"domain.events",
		"notifications.meetup_events",
		"meetup.created",
	)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ subscriber: %v", err)
	}
	defer subscriber.Close()

	repo := infrastructure.NewPostgresRepo(db)

	worker := infrastructure.NewConsumerWorker(repo, msgsChan)
	go worker.Start(context.Background())

	mux := http.NewServeMux()

	shared.RegisterHealthRoutes(mux, db.DB, "notifications", CommitHash)

	getNotificationsHandler := get_notifications.NewHandler(repo)
	api.RegisterRoutes(mux, getNotificationsHandler)

	log.Printf("Starting Notification Service on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
