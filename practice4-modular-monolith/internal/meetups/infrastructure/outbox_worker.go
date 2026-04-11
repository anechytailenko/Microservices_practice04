package infrastructure

import (
	"context"
	"database/sql"
	"log"
	"time"
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

	log.Println("[Outbox Worker] Started processing background events...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Outbox Worker] Shutting down...")
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
        WHERE processed_at IS NULL 
        ORDER BY created_at ASC 
        LIMIT 50
        FOR UPDATE SKIP LOCKED`

	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[Outbox Worker] Error querying database: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id, eventType string
		var payload []byte

		if err := rows.Scan(&id, &eventType, &payload); err != nil {
			log.Printf("[Outbox Worker] Error scanning row: %v\n", err)
			continue
		}

		err = w.publisher.Publish(ctx, eventType, payload)
		if err != nil {
			log.Printf("[Outbox Worker] Failed to publish EventID %s: %v. Will retry later.", id, err)
			continue
		}

		updateQuery := `UPDATE outbox_events SET processed_at = NOW() WHERE id = $1`
		_, err = w.db.ExecContext(ctx, updateQuery, id)

		if err != nil {
			log.Printf("[Outbox Worker] Failed to update processed_at for %s: %v\n", id, err)
		} else {
			log.Printf("[Outbox Worker] <-- SUCCESSFULLY PUBLISHED AND PROCESSED: %s (Type: %s)\n", id, eventType)
		}
	}
}
