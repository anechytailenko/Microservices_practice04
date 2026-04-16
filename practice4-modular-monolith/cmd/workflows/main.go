package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/database"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
	workflow_api "github.com/anechytailenko/Microservices_practice04/internal/workflows/api"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/features/create_workflow"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/features/get_workflow"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/migrations"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var CommitHash string
	if envHash := os.Getenv("COMMIT_HASH"); envHash != "" {
		CommitHash = envHash
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("WORKFLOW_DB_URL")
	if dbURL == "" {
		logger.Fatal(context.Background(), "WORKFLOW_DB_URL is not set")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Fatal(context.Background(), "RABBITMQ_URL is not set")
	}

	exchangeName := os.Getenv("DOMAIN_EXCHANGE")
	if exchangeName == "" {
		exchangeName = "domain.events"
	}

	queueName := os.Getenv("WORKFLOW_QUEUE")
	if queueName == "" {
		queueName = "workflow.saga_events"
	}

	dlxName := os.Getenv("EVENTS_DLX")
	if dlxName == "" {
		dlxName = "events.dlx"
	}

	dlqName := os.Getenv("WORKFLOW_DLQ")
	if dlqName == "" {
		dlqName = "workflow.dlq"
	}

	consumerName := os.Getenv("CONSUMER_NAME")
	if consumerName == "" {
		consumerName = "workflow_service"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to open database: %v", err)
	}
	defer db.Close()

	// process of migration that will run as a separate Job in kubernetes
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		logger.Println(context.Background(), "Starting migration process for workflows...")

		if err := database.RunMigrations(db, migrations.FS, "workflows"); err != nil {
			logger.Fatalf(context.Background(), "Fatal error running workflow migrations: %v", err)
		}

		logger.Println(context.Background(), "Migrations finished successfully. Exiting.")

		db.Close()
		os.Exit(0)
	}

	publisher, err := rabbitmq.NewPublisher(rabbitMQURL, exchangeName)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

	if err := database.RunMigrations(db, migrations.FS, "workflows"); err != nil {
		logger.Fatalf(context.Background(), "Fatal error running workflow migrations: %v", err)
	}

	bindingKeys := []string{"events.#"}

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
		logger.Fatalf(context.Background(), "Failed to initialize RabbitMQ subscriber: %v", err)
	}
	defer subscriber.Close()

	repo := infrastructure.NewPostgresRepo(db)

	outboxWorker := infrastructure.NewOutboxWorker(db, publisher)
	go outboxWorker.Start(context.Background(), 2*time.Second)

	consumerWorker := infrastructure.NewConsumerWorker(db, repo, msgsChan)
	go consumerWorker.Start(context.Background())

	timeoutWorker := infrastructure.NewTimeoutWorker(db, repo, 30*time.Second)
	go timeoutWorker.Start(context.Background(), 10*time.Second)

	createHandler := create_workflow.NewHandler(repo)
	getHandler := get_workflow.NewHandler(repo)

	mux := http.NewServeMux()

	shared.RegisterHealthRoutes(mux, db, "workflow", CommitHash)
	workflow_api.RegisterRoutes(mux, createHandler, getHandler)

	finalHandler := shared.ContextWithCorrelationID(mux)

	logger.Printf(context.Background(), "Starting Workflow Service on port %s...", port)
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		logger.Fatalf(context.Background(), "Server failed: %v", err)
	}
}
