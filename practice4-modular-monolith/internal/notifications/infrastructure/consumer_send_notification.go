package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) handleSendNotification(ctx context.Context, d amqp.Delivery) {
	var cmd commands.SendNotification
	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		logger.Println(ctx, "[Notifications Consumer] Invalid JSON for SendNotification: %v", err)
		d.Nack(false, false)
		return
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		d.Nack(false, true)
		return
	}
	defer tx.Rollback()

	inboxID := fmt.Sprintf("%s-%s", cmd.WorkflowID, commands.SendNotificationType)
	if !w.checkInbox(ctx, tx, inboxID) {
		logger.Println(ctx, "[Notifications Consumer] Saga Notification %s already processed. Skipping.", inboxID)
		d.Ack(false)
		return
	}

	eventID := uuid.New().String()
	payloadJSON, _ := json.Marshal(cmd)

	err = w.repo.SaveNotificationTx(ctx, tx, eventID, cmd.CorrelationID, cmd.UserID, payloadJSON)
	if err != nil {
		logger.Println(ctx, "[Notifications Consumer] Failed to save Saga notification: %v", err)
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		d.Nack(false, true)
		return
	}

	logger.Println(ctx, "[Notifications Consumer] Sent %s notification to User %s. Msg: %s", cmd.Status, cmd.UserID, cmd.Message)
	d.Ack(false)
}
