package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) processUpdateSchedule(ctx context.Context, d amqp.Delivery) {
	var cmd commands.UpdateSchedule
	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		log.Printf("[Users Consumer] Invalid JSON for UpdateSchedule: %v", err)
		d.Nack(false, false)
		return
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		d.Nack(false, true)
		return
	}
	defer tx.Rollback()

	inboxEventID := fmt.Sprintf("%s-%s", cmd.WorkflowID, commands.UpdateScheduleType)
	if !w.checkInbox(ctx, tx, inboxEventID) {
		log.Printf("[Users Consumer] Command %s already processed. Skipping.", inboxEventID)
		d.Ack(false)
		return
	}

	user, err := w.repo.GetByIDTx(ctx, tx, cmd.UserID)
	if err != nil {
		w.handleUpdateErrorFlow(ctx, err, tx, d, cmd)
		return
	}

	user.AddMeetup(cmd.MeetupID)

	err = w.repo.UpdateUserMeetupsTx(ctx, tx, user)
	if err != nil {
		log.Printf("[Users Consumer] Failed to update user meetups in DB: %v", err)
		d.Nack(false, true)
		return
	}

	replyEvent := events.ScheduleUpdated{
		WorkflowID:    cmd.WorkflowID,
		UserID:        cmd.UserID,
		MeetupID:      cmd.MeetupID,
		CorrelationID: cmd.CorrelationID,
	}

	payload, err := json.Marshal(replyEvent)
	if err != nil {
		d.Nack(false, true)
		return
	}

	err = w.repo.SaveOutboxEventTx(ctx, tx, events.ScheduleUpdatedEventType, payload)
	if err != nil {
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		d.Nack(false, true)
		return
	}

	log.Printf("[Users Consumer] Successfully updated schedule for User %s (Workflow %s)", cmd.UserID, cmd.WorkflowID)
	d.Ack(false)
}

func (w *ConsumerWorker) handleUpdateErrorFlow(ctx context.Context, err error, tx *sql.Tx, d amqp.Delivery, cmd commands.UpdateSchedule) {
	var sharedErr shared.Error

	if errors.As(err, &sharedErr) && sharedErr.Type == shared.ErrorTypeNotFound {
		log.Printf("[Users Consumer] User not found. Failing workflow %s.", cmd.WorkflowID)

		replyEvent := events.ScheduleUpdateFailed{
			WorkflowID:    cmd.WorkflowID,
			UserID:        cmd.UserID,
			MeetupID:      cmd.MeetupID,
			Reason:        fmt.Sprintf("User %s does not exist", cmd.UserID),
			CorrelationID: cmd.CorrelationID,
		}

		payload, _ := json.Marshal(replyEvent)

		saveErr := w.repo.SaveOutboxEventTx(ctx, tx, events.ScheduleUpdateFailedEventType, payload)
		if saveErr != nil {
			d.Nack(false, true)
			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			d.Nack(false, true)
			return
		}

		d.Ack(false)
		return
	}

	log.Printf("[Users Consumer] Internal error for Workflow %s: %v", cmd.WorkflowID, err)
	d.Nack(false, true)
}
