package infrastructure

import (
	"context"
	"encoding/json"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) handleMeetupCreated(ctx context.Context, d amqp.Delivery) {
	var evt events.MeetupCreatedEvent
	if err := json.Unmarshal(d.Body, &evt); err != nil {
		logger.Println(ctx, "[Notifications Consumer] Invalid JSON for MeetupCreated: %v", err)
		d.Nack(false, false)
		return
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		d.Nack(false, true)
		return
	}
	defer tx.Rollback()

	if !w.checkInbox(ctx, tx, evt.EventID) {
		logger.Println(ctx, "[Notifications Consumer] Legacy Event %s already processed. Skipping.", evt.EventID)
		d.Ack(false)
		return
	}

	err = w.repo.SaveNotificationTx(ctx, tx, evt.EventID, evt.CorrelationID, evt.OwnerUserID, d.Body)
	if err != nil {
		logger.Println(ctx, "[Notifications Consumer] Failed to save MeetupCreated notification: %v", err)
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		d.Nack(false, true)
		return
	}

	logger.Println(ctx, "[Notifications Consumer] Successfully processed MeetupCreated Event %s", evt.EventID)
	d.Ack(false)
}
