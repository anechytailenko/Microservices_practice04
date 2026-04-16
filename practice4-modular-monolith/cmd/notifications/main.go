package main

import (
	"context"
	"net/http"
	"os"

	"github.com/anechytailenko/Microservices_practice04/internal/notifications/api"
	"github.com/anechytailenko/Microservices_practice04/internal/notifications/features/get_notifications"
	"github.com/anechytailenko/Microservices_practice04/internal/notifications/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/database"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
	"github.com/anechytailenko/Microservices_practice04/migrations"

	"github.com/jmoiron/sqlx"
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

	dbURL := os.Getenv("NOTIFICATIONS_DB_URL")
	if dbURL == "" {
		logger.Fatal(context.Background(), "NOTIFICATIONS_DB_URL is not set")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		logger.Fatal(context.Background(), "RABBITMQ_URL is not set")
	}

	exchangeName := os.Getenv("DOMAIN_EXCHANGE")
	if exchangeName == "" {
		exchangeName = "domain.events"
	}

	queueName := os.Getenv("NOTIFICATIONS_QUEUE")
	if queueName == "" {
		queueName = "notifications.meetup_events"
	}

	dlxName := os.Getenv("EVENTS_DLX")
	if dlxName == "" {
		dlxName = "events.dlx"
	}

	dlqName := os.Getenv("NOTIFICATIONS_DLQ")
	if dlqName == "" {
		dlqName = "notifications.dlq"
	}

	consumerName := os.Getenv("CONSUMER_NAME")
	if consumerName == "" {
		consumerName = "notification_service"
	}

	db, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		logger.Fatalf(context.Background(), "Failed to connect to database: %v", err)
	}
	defer db.Close()

	// process of migration that will run as a separate Job in kubernetes
	if len(os.Args) > 1 && os.Args[1] == "migrate" {

		logger.Println(context.Background(), "Starting migration process for notifications...")

		if err := database.RunMigrations(db.DB, migrations.FS, "notifications"); err != nil {
			logger.Fatalf(context.Background(), "Fatal error running notifications migrations: %v", err)
		}

		logger.Println(context.Background(), "Migrations finished successfully. Exiting.")

		db.Close()
		os.Exit(0)
	}

	if err := database.RunMigrations(db.DB, migrations.FS, "notifications"); err != nil {
		logger.Fatalf(context.Background(), "Fatal error running notifications migrations: %v", err)
	}

	bindingKeys := []string{
		"commands.notifications.#",
		"events.meetups.created",
	}

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

	worker := infrastructure.NewConsumerWorker(db.DB, repo, msgsChan)
	go worker.Start(context.Background())

	mux := http.NewServeMux()

	shared.RegisterHealthRoutes(mux, db.DB, "notifications", CommitHash)

	getNotificationsHandler := get_notifications.NewHandler(repo)
	api.RegisterRoutes(mux, getNotificationsHandler)

	finalHandler := shared.ContextWithCorrelationID(mux)

	logger.Printf(context.Background(), "Starting Notification Service on port %s...", port)
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		logger.Fatalf(context.Background(), "Server failed: %v", err)
	}
}
