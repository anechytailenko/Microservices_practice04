package infrastructure

import (
	"context"
	"database/sql"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
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
	logger.Println(ctx, "[Workflow Consumer] Started listening for saga events...")

	for {
		select {
		case <-ctx.Done():
			logger.Println(ctx, "[Workflow Consumer] Shutting down gracefully...")
			return
		case d, ok := <-w.msgs:
			if !ok {
				logger.Println(ctx, "[Workflow Consumer] Message channel closed")
				return
			}

			var corrID string
			if d.Headers != nil {
				if id, ok := d.Headers["X-Correlation-Id"].(string); ok {
					corrID = id
				}
			}

			msgCtx := ctxutil.WithCorrelationID(ctx, corrID)
			w.processEvent(msgCtx, d)
		}
	}
}

func (w *ConsumerWorker) processEvent(ctx context.Context, d amqp.Delivery) {
	switch d.RoutingKey {
	case events.SeatReservedEventType:
		w.handleSeatReserved(ctx, d)
	case events.SeatReservationFailedEventType:
		w.handleSeatReservationFailed(ctx, d)
	case events.ScheduleUpdatedEventType:
		w.handleScheduleUpdated(ctx, d)
	case events.ScheduleUpdateFailedEventType:
		w.handleScheduleUpdateFailed(ctx, d)
	case events.SeatReservationCanceledEventType:
		w.handleSeatCanceled(ctx, d)
	default:
		logger.Printf(ctx, "[Workflow Consumer] Unknown event type: %s. Dropping.", d.RoutingKey)
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
		logger.Printf(ctx, "[Workflow Consumer] Inbox DB Error: %v", err)
		return false
	}
	count, _ := res.RowsAffected()
	return count > 0
}
