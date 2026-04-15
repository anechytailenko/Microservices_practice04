package infrastructure

import (
	"context"
	"database/sql"
	"log"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerWorker struct {
	db       *sql.DB
	repo     *PostgresRepo
	messages <-chan amqp.Delivery
}

func NewConsumerWorker(db *sql.DB, messages <-chan amqp.Delivery) *ConsumerWorker {
	return &ConsumerWorker{
		db:       db,
		repo:     NewPostgresRepo(db),
		messages: messages,
	}
}

func (w *ConsumerWorker) Start(ctx context.Context) {
	log.Println("[Meetups Consumer] Started listening for commands...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Meetups Consumer] Shutting down gracefully...")
			return
		case d, ok := <-w.messages:
			if !ok {
				log.Println("[Meetups Consumer] Message channel closed")
				return
			}
			w.processEvent(ctx, d)
		}
	}
}

func (w *ConsumerWorker) processEvent(ctx context.Context, d amqp.Delivery) {
	switch d.RoutingKey {
	case commands.ReserveSeatType:
		w.processReserveSeat(ctx, d)
	case commands.CancelSeatReservationType:
		w.processCancelSeat(ctx, d)
	default:
		log.Printf("[Meetups Consumer] Unknown routing key: %s. Dropping message.", d.RoutingKey)
		d.Nack(false, false)
	}
}

// inbox patter - that guarantee idempodency
func (w *ConsumerWorker) checkInbox(ctx context.Context, tx *sql.Tx, eventID string) bool {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO inbox_events (event_id) 
		VALUES ($1) 
		ON CONFLICT (event_id) DO NOTHING`,
		eventID,
	)
	if err != nil {
		log.Printf("[Meetups Consumer] Inbox DB Error for Event %s: %v", eventID, err)
		return false
	}
	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0
}
