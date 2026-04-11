package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConsumerWorker struct {
	repo     *PostgresRepo
	messages <-chan amqp.Delivery
}

func NewConsumerWorker(repo *PostgresRepo, msgs <-chan amqp.Delivery) *ConsumerWorker {
	return &ConsumerWorker{
		repo:     repo,
		messages: msgs,
	}
}

func (w *ConsumerWorker) Start(ctx context.Context) {
	log.Println("[Notification Consumer] Started listening for events...")

	for {
		select {
		case <-ctx.Done():
			log.Println("[Notification Consumer] Shutting down gracefully...")
			return

		case d, ok := <-w.messages:
			if !ok {
				log.Println("[Notification Consumer] Message channel closed")
				return
			}

			w.processEvent(ctx, d)
		}
	}
}

func (w *ConsumerWorker) processEvent(ctx context.Context, d amqp.Delivery) {
	var partialEvent struct {
		EventID       string `json:"eventId"`
		CorrelationID string `json:"correlationId"`
		OwnerUserID   string `json:"ownerUserId"`
	}

	if err := json.Unmarshal(d.Body, &partialEvent); err != nil {
		log.Printf("[Notification Consumer] Invalid JSON payload: %v. NACKing (sending to DLQ)...", err)
		d.Nack(false, false)
		return
	}

	err := w.repo.SaveNotification(
		ctx,
		partialEvent.EventID,
		partialEvent.CorrelationID,
		partialEvent.OwnerUserID,
		d.Body,
	)

	if err != nil {
		log.Printf("[Notification Consumer] DB Error for Event %s: %v. Requeueing...", partialEvent.EventID, err)
		d.Nack(false, true)
	} else {
		log.Printf("[Notification Consumer] Successfully processed Event %s", partialEvent.EventID)
		d.Ack(false)
	}
}
