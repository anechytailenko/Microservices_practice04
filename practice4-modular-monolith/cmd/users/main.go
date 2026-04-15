package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

<<<<<<< Updated upstream
	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
=======
<<<<<<< Updated upstream
<<<<<<< Updated upstream
=======
=======
>>>>>>> Stashed changes
	"github.com/anechytailenko/Microservices_practice04/internal/shared/database"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
>>>>>>> Stashed changes
>>>>>>> Stashed changes
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
	users_api "github.com/anechytailenko/Microservices_practice04/internal/users/api"
	"github.com/anechytailenko/Microservices_practice04/internal/users/features/create_user"
	"github.com/anechytailenko/Microservices_practice04/internal/users/features/get_user"
	"github.com/anechytailenko/Microservices_practice04/internal/users/infrastructure"
	"github.com/anechytailenko/Microservices_practice04/migrations"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// this var will be rewritten by compiler when we would build

func main() {
	var CommitHash string
	if envHash := os.Getenv("COMMIT_HASH"); envHash != "" {
		CommitHash = envHash
	}

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

<<<<<<< Updated upstream
=======
<<<<<<< Updated upstream
<<<<<<< Updated upstream
=======
=======
>>>>>>> Stashed changes
	// proccess of migration that will run as the seperate Job in kubernetes
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		log.Println("Starting migration process for users...")
		if err := database.RunMigrations(db, migrations.FS, "users"); err != nil {
			log.Fatalf("Fatal error running users migrations: %v", err)
		}
		log.Println("Migrations finished successfully. Exiting.")
		db.Close()
		os.Exit(0)
	}

>>>>>>> Stashed changes
	publisher, err := rabbitmq.NewPublisher(rabbitMQURL, exchangeName)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

<<<<<<< Updated upstream
	bindingKeys := []string{"commands.users.*"}
=======
	if err := database.RunMigrations(db, migrations.FS, "users"); err != nil {
		log.Fatalf("Fatal error running users migrations: %v", err)
	}

	bindingKeys := []string{"commands.users.#"}
>>>>>>> Stashed changes

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

<<<<<<< Updated upstream
=======
>>>>>>> Stashed changes
>>>>>>> Stashed changes
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
