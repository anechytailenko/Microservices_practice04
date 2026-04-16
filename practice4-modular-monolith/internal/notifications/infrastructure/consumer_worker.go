package infrastructure

import (
	"context"
	"database/sql"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerWorker struct {
	db       *sql.DB
	repo     *PostgresRepo
	messages <-chan amqp.Delivery
}

func NewConsumerWorker(db *sql.DB, repo *PostgresRepo, msgs <-chan amqp.Delivery) *ConsumerWorker {
	return &ConsumerWorker{
		db:       db,
		repo:     repo,
		messages: msgs,
	}
}

func (w *ConsumerWorker) Start(ctx context.Context) {
	logger.Println(ctx, "[Notifications Consumer] Started listening for events & commands...")

	for {
		select {
		case <-ctx.Done():
			logger.Println(ctx, "[Notifications Consumer] Shutting down gracefully...")
			return
		case d, ok := <-w.messages:
			if !ok {
				logger.Println(ctx, "[Notifications Consumer] Message channel closed")
				return
			}
			corrID := d.Headers["X-Correlation-Id"].(string)
			ctx := ctxutil.WithCorrelationID(context.Background(), corrID)
			w.processEvent(ctx, d)
		}
	}
}

func (w *ConsumerWorker) processEvent(ctx context.Context, d amqp.Delivery) {
	switch d.RoutingKey {
	case commands.SendNotificationType:
		w.handleSendNotification(ctx, d)
	case events.MeetupCreatedEventType:
		w.handleMeetupCreated(ctx, d)
	default:
		logger.Println(ctx, "[Notifications Consumer] Unknown event type: %s. Dropping.", d.RoutingKey)
		d.Ack(false)
	}
}

func (w *ConsumerWorker) checkInbox(ctx context.Context, tx *sql.Tx, id string) bool {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO inbox_events (event_id) 
		VALUES ($1) 
		ON CONFLICT (event_id) DO NOTHING`,
		id,
	)
	if err != nil {
		logger.Println(ctx, "[Notifications Consumer] Inbox DB Error: %v", err)
		return false
	}
	count, _ := res.RowsAffected()
	return count > 0
}
