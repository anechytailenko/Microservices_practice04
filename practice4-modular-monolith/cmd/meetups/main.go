package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	meetups_api "github.com/anechytailenko/Microservices_practice04/internal/meetups/api"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/database"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
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

	dbURL := os.Getenv("MEETUPS_DB_URL")
	if dbURL == "" {
		logger.Fatal(context.Background(), "MEETUPS_DB_URL is not set")
	}

	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		logger.Fatal(context.Background(), "USERS_SERVICE_URL is not set")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Fatal(context.Background(), "RABBITMQ_URL is not set")
	}

	exchangeName := os.Getenv("DOMAIN_EXCHANGE")
	if exchangeName == "" {
		exchangeName = "domain.events"
	}

	queueName := os.Getenv("MEETUPS_QUEUE")
	if queueName == "" {
		queueName = "meetups.workflow_commands"
	}

	dlxName := os.Getenv("EVENTS_DLX")
	if dlxName == "" {
		dlxName = "events.dlx"
	}

	dlqName := os.Getenv("MEETUPS_DLQ")
	if dlqName == "" {
		dlqName = "meetups.dlq"
	}

	consumerName := os.Getenv("CONSUMER_NAME")
	if consumerName == "" {
		consumerName = "meetups_service"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to open database: %v", err)
	}
	defer db.Close()

	// process of migration that will run as a separate Job in kubernetes
	if len(os.Args) > 1 && os.Args[1] == "migrate" {

		logger.Println(context.Background(), "Starting migration process for meetups...")

		if err := database.RunMigrations(db, migrations.FS, "meetups"); err != nil {
			logger.Fatalf(context.Background(), "Fatal error running meetups migrations: %v", err)
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

	bindingKeys := []string{"commands.meetups.#"}

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

	outboxWorker := infrastructure.NewOutboxWorker(db, publisher)
	go outboxWorker.Start(context.Background(), 2*time.Second)

	consumerWorker := infrastructure.NewConsumerWorker(db, msgsChan)
	go consumerWorker.Start(context.Background())

	repo := infrastructure.NewPostgresRepo(db)
	usersClient := infrastructure.NewUsersClient(usersServiceURL)

	createHandler := create_meetup.NewHandler(repo, usersClient)
	changeStatusHandler := change_status.NewHandler(repo)
	getMeetupHandler := get_meetup.NewHandler(repo)

	mux := http.NewServeMux()

	shared.RegisterHealthRoutes(mux, db, "meetups", CommitHash)
	meetups_api.RegisterRoutes(mux, createHandler, changeStatusHandler, getMeetupHandler)

	finalHandler := shared.ContextWithCorrelationID(mux)

	logger.Printf(context.Background(), "Starting Meetups Service on port %s...", port)
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		logger.Fatalf(context.Background(), "Server failed: %v", err)
	}
}
