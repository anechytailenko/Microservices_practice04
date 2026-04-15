package infrastructure

import (
	"context"
	"database/sql"
	"log"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerWorker struct {
	db   *sql.DB
	repo *PostgresRepo
	msgs <-chan amqp.Delivery
}

func NewConsumerWorker(db *sql.DB, repo *PostgresRepo, msgs <-chan amqp.Delivery) *ConsumerWorker {
	return &ConsumerWorker{
		db:   db,
		repo: repo,
		msgs: msgs,
	}
}

func (w *ConsumerWorker) Start(ctx context.Context) {
	log.Println("[Users Consumer] Started listening for commands...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Users Consumer] Shutting down gracefully...")
			return
		case d, ok := <-w.msgs:
			if !ok {
				log.Println("[Users Consumer] Message channel closed")
				return
			}
			w.processCommand(ctx, d)
		}
	}
}

func (w *ConsumerWorker) processCommand(ctx context.Context, d amqp.Delivery) {
	switch d.RoutingKey {
	case commands.UpdateScheduleType:
		w.processUpdateSchedule(ctx, d)
	default:
		log.Printf("[Users Consumer] Unknown routing key: %s. Dropping.", d.RoutingKey)
		d.Ack(false)
	}
}

// inbox pattern
func (w *ConsumerWorker) checkInbox(ctx context.Context, tx *sql.Tx, eventID string) bool {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO inbox_events (event_id) 
		VALUES ($1) 
		ON CONFLICT (event_id) DO NOTHING`,
		eventID,
	)
	if err != nil {
		log.Printf("[Users Consumer] Inbox DB Error: %v", err)
		return false
	}
	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0
}
