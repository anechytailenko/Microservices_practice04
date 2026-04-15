package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
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

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
	}

	exchangeName := os.Getenv("DOMAIN_EXCHANGE")
	if exchangeName == "" {
		exchangeName = "domain.events"
	}

	queueName := os.Getenv("USERS_QUEUE")
	if queueName == "" {
		queueName = "users.workflow_commands"
	}

	dlxName := os.Getenv("EVENTS_DLX")
	if dlxName == "" {
		dlxName = "events.dlx"
	}

	dlqName := os.Getenv("USERS_DLQ")
	if dlqName == "" {
		dlqName = "users.dlq"
	}

	consumerName := os.Getenv("CONSUMER_NAME")
	if consumerName == "" {
		consumerName = "users_service"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	publisher, err := rabbitmq.NewPublisher(rabbitMQURL, exchangeName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

	bindingKeys := []string{"commands.users.*"}

	subscriber, msgsChan, err := rabbitmq.NewSubscriber(
		rabbitMQURL,
		exchangeName,
		queueName,
		bindingKeys,
		dlxName,
		dlqName,
		consumerName,
	)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ subscriber: %v", err)
	}
	defer subscriber.Close()

	repo := infrastructure.NewPostgresRepo(db)

	outboxWorker := infrastructure.NewOutboxWorker(db, publisher)
	go outboxWorker.Start(context.Background(), 2*time.Second)

	consumerWorker := infrastructure.NewConsumerWorker(db, repo, msgsChan)
	go consumerWorker.Start(context.Background())

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
