package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	meetups_api "github.com/anechytailenko/Microservices_practice04/internal/meetups/api"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/database"
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
		log.Fatal("MEETUPS_DB_URL is not set")
	}

	usersServiceURL := os.Getenv("USERS_SERVICE_URL")
	if usersServiceURL == "" {
		log.Fatal("USERS_SERVICE_URL is not set")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL is not set")
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
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// proccess of migration that will run as the seperate Job in kubernetes
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Starting  migration process for meetups...")
		if err := database.RunMigrations(db, migrations.FS, "meetups"); err != nil {
			log.Fatalf("Fatal error running meetups migrations: %v", err)
		}
		log.Println("Migrations finished successfully. Exiting.")
		db.Close()
		os.Exit(0)
	}

	publisher, err := rabbitmq.NewPublisher(rabbitMQURL, exchangeName)

	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
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
		log.Fatalf("Failed to initialize RabbitMQ subscriber: %v", err)
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

	log.Printf("Starting Meetups Service on port %s...", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
