package infrastructure

import (
	"context"
	"database/sql"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
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

	logger.Println(ctx, "[Users Consumer] Started listening for commands...")

	for {
		select {
		case <-ctx.Done():
			logger.Println(ctx, "[Users Consumer] Shutting down gracefully...")
			return
		case d, ok := <-w.msgs:
			if !ok {
				logger.Println(ctx, "[Users Consumer] Message channel closed")
				return
			}

			var corrID string
			if d.Headers != nil {
				if id, ok := d.Headers["X-Correlation-Id"].(string); ok {
					corrID = id
				}
			}

			msgCtx := ctxutil.WithCorrelationID(ctx, corrID)
			w.processCommand(msgCtx, d)
		}
	}
}

func (w *ConsumerWorker) processCommand(ctx context.Context, d amqp.Delivery) {
	switch d.RoutingKey {
	case commands.UpdateScheduleType:
		w.processUpdateSchedule(ctx, d)
	default:
		logger.Printf(ctx, "[Users Consumer] Unknown routing key: %s. Dropping.", d.RoutingKey)
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
		logger.Printf(ctx, "[Users Consumer] Inbox DB Error: %v", err)
		return false
	}
	rowsAffected, _ := res.RowsAffected()
	return rowsAffected > 0
}
