package infrastructure

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/rabbitmq"
	"github.com/lib/pq"
)

type OutboxWorker struct {
	db        *sql.DB
	publisher *rabbitmq.Publisher
}

func NewOutboxWorker(db *sql.DB, publisher *rabbitmq.Publisher) *OutboxWorker {
	return &OutboxWorker{
		db:        db,
		publisher: publisher,
	}
}

func (w *OutboxWorker) Start(ctx context.Context, interval time.Duration) {
	log.Printf("[Users Outbox Worker] Started processing events every %v...", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[Users Outbox Worker] Shutting down gracefully...")
			return
		case <-ticker.C:
			w.processOutbox(ctx)
		}
	}
}

func (w *OutboxWorker) processOutbox(ctx context.Context) {
	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[Users Outbox Worker] Failed to begin tx: %v", err)
		return
	}
	defer tx.Rollback()

	query := `
		SELECT id, event_type, payload 
		FROM outbox_events 
		WHERE processed = false 
		ORDER BY created_at ASC 
		LIMIT 50 
		FOR UPDATE SKIP LOCKED`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		log.Printf("[Users Outbox Worker] DB Error fetching events: %v", err)
		return
	}
	defer rows.Close()

	type outboxEvent struct {
		ID        string
		EventType string
		Payload   []byte
	}

	var eventsToProcess []outboxEvent
	for rows.Next() {
		var evt outboxEvent
		if err := rows.Scan(&evt.ID, &evt.EventType, &evt.Payload); err != nil {
			continue
		}
		eventsToProcess = append(eventsToProcess, evt)
	}
	rows.Close()

	if len(eventsToProcess) == 0 {
		tx.Commit()
		return
	}

	var processedIDs []string

	for _, evt := range eventsToProcess {
		err := w.publisher.Publish(ctx, evt.EventType, evt.Payload)
		if err != nil {
			log.Printf("[Users Outbox Worker] Failed to publish event %s: %v", evt.ID, err)
			continue
		}
		processedIDs = append(processedIDs, evt.ID)
	}

	if len(processedIDs) == 0 {
		tx.Commit()
		return
	}

	updateQuery := `
		UPDATE outbox_events 
		SET processed = true 
		WHERE id = ANY($1)`

	_, err = tx.ExecContext(ctx, updateQuery, pq.Array(processedIDs))
	if err != nil {
		log.Printf("[Users Outbox Worker] Failed to mark events as processed: %v", err)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[Users Outbox Worker] Failed to commit outbox processing: %v", err)
		return
	}

	log.Printf("[Users Outbox Worker] Successfully published %d events to RabbitMQ", len(processedIDs))
}
