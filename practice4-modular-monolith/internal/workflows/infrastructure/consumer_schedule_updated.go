package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) handleScheduleUpdated(ctx context.Context, d amqp.Delivery) {
	var evt events.ScheduleUpdated
	if err := json.Unmarshal(d.Body, &evt); err != nil {
		d.Nack(false, false)
		return
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		d.Nack(false, true)
		return
	}
	defer tx.Rollback()

	inboxID := fmt.Sprintf("%s-%s", evt.WorkflowID, events.ScheduleUpdatedEventType)
	if !w.checkInbox(ctx, tx, inboxID) {
		d.Ack(false)
		return
	}

	wf, err := w.repo.GetByIDTx(ctx, tx, domain.WorkflowID(evt.WorkflowID))
	if err != nil {
		d.Nack(false, true)
		return
	}

	err = wf.ChangeState(domain.StateCompleted)
	if err != nil {
		var sharedErr shared.Error
		if errors.As(err, &sharedErr) && sharedErr.Type == shared.ErrorTypeConflict {
			log.Printf("[Workflow Consumer] Dropping late event: %v", err.Error())
			d.Ack(false)
			return
		}
		d.Nack(false, true)
		return
	}

	err = w.repo.UpdateTx(ctx, tx, wf)
	if err != nil {
		d.Nack(false, true)
		return
	}

	notifyCmd := commands.SendNotification{
		WorkflowID:    evt.WorkflowID,
		UserID:        evt.UserID,
		MeetupID:      evt.MeetupID,
		Status:        "success",
		Message:       "Hi, here is the invitation to meet up...",
		CorrelationID: evt.CorrelationID,
	}

	notifyPayload, _ := json.Marshal(notifyCmd)
	err = w.repo.SaveOutboxEventTx(ctx, tx, commands.SendNotificationType, notifyPayload)
	if err != nil {
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		d.Nack(false, true)
		return
	}

	log.Printf("[Workflow] Saga %s: SUCCESSFULLY COMPLETED ", evt.WorkflowID)
	d.Ack(false)
}
