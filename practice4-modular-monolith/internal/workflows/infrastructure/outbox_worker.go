package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
)

type EventPublisher interface {
	Publish(ctx context.Context, routingKey string, payload []byte) error
}

type OutboxWorker struct {
	db        *sql.DB
	publisher EventPublisher
}

func NewOutboxWorker(db *sql.DB, publisher EventPublisher) *OutboxWorker {
	return &OutboxWorker{
		db:        db,
		publisher: publisher,
	}
}

func (w *OutboxWorker) Start(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	logger.Println(ctx, "[Workflow Outbox Worker] Started processing background events...")

	for {
		select {
		case <-ctx.Done():
			logger.Println(ctx, "[Workflow Outbox Worker] Shutting down...")
			return
		case <-ticker.C:
			w.processOutbox(ctx)
		}
	}
}

func (w *OutboxWorker) processOutbox(ctx context.Context) {
	query := `
        SELECT id, event_type, payload 
        FROM outbox_events 
        WHERE processed = false 
        ORDER BY created_at ASC 
        LIMIT 50
        FOR UPDATE SKIP LOCKED`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		logger.Printf(ctx, "[Workflow Outbox Worker] Error querying database: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, eventType string
		var payload []byte

		if err := rows.Scan(&id, &eventType, &payload); err != nil {
			logger.Printf(ctx, "[Workflow Outbox Worker] Error scanning row: %v", err)
			continue
		}

		err = w.publisher.Publish(ctx, eventType, payload)
		if err != nil {
			logger.Printf(ctx, "[Workflow Outbox Worker] Failed to publish EventID %s: %v. Will retry next tick.", id, err)
			continue
		}

		updateQuery := `UPDATE outbox_events SET processed = true WHERE id = $1`
		_, err = w.db.ExecContext(ctx, updateQuery, id)

		if err != nil {
			logger.Printf(ctx, "[Workflow Outbox Worker] CRITICAL: Failed to mark %s as processed: %v", id, err)
		} else {
			logger.Printf(ctx, "[Workflow Outbox Worker] --> PUBLISHED: %s (Key: %s)", id, eventType)
		}
	}
}
