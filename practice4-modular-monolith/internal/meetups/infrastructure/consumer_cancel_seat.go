package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) processCancelSeat(ctx context.Context, d amqp.Delivery) {
	var cmd commands.CancelSeatReservation

	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Invalid JSON for CancelSeat: %v", err)
		d.Nack(false, false)
		return
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Failed to begin transaction: %v", err)
		d.Nack(false, true)
		return
	}
	defer tx.Rollback()

	inboxEventID := fmt.Sprintf("%s-%s", cmd.WorkflowID, commands.CancelSeatReservationType)
	if !w.checkInbox(ctx, tx, inboxEventID) {
		logger.Printf(ctx, "[Meetups Consumer] Cancel event %s already processed. Skipping.", inboxEventID)
		d.Ack(false)
		return
	}

	meetup, err := w.repo.GetByIDTx(ctx, tx, cmd.MeetupID)
	if err != nil {
		w.handleCancelErrorFlow(ctx, err, tx, d, cmd)
		return
	}

	meetup.RemoveGuest(cmd.UserID)

	err = w.repo.UpdateGuestsTx(ctx, tx, meetup)
	if err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Failed to update DB for cancellation: %v", err)
		d.Nack(false, true)
		return
	}

	replyEvent := events.SeatReservationCanceled{
		WorkflowID:    cmd.WorkflowID,
		MeetupID:      cmd.MeetupID,
		UserID:        cmd.UserID,
		CorrelationID: cmd.CorrelationID,
	}

	payload, err := json.Marshal(replyEvent)
	if err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Failed to marshal reply event: %v", err)
		d.Nack(false, true)
		return
	}

	outboxEventID := uuid.New().String()
	err = w.repo.SaveOutboxEventTx(ctx, tx, outboxEventID, events.SeatReservationCanceledEventType, payload)
	if err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Failed to save outbox cancel event: %v", err)
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Commit failed for cancellation: %v", err)
		d.Nack(false, true)
		return
	}

	logger.Printf(ctx, "[Meetups Consumer] Successfully canceled seat reservation for Workflow %s", cmd.WorkflowID)
	d.Ack(false)
}

func (w *ConsumerWorker) handleCancelErrorFlow(ctx context.Context, err error, tx *sql.Tx, d amqp.Delivery, cmd commands.CancelSeatReservation) {
	var sharedErr shared.Error

	if errors.As(err, &sharedErr) && sharedErr.Type == shared.ErrorTypeNotFound {
		logger.Printf(ctx, "[Meetups Consumer] Meetup not found during cancellation for Workflow %s. Assuming success.", cmd.WorkflowID)

		replyEvent := events.SeatReservationCanceled{
			WorkflowID:    cmd.WorkflowID,
			MeetupID:      cmd.MeetupID,
			UserID:        cmd.UserID,
			CorrelationID: cmd.CorrelationID,
		}

		payload, marshalErr := json.Marshal(replyEvent)
		if marshalErr != nil {
			logger.Printf(ctx, "[Meetups Consumer] Failed to marshal compensation reply: %v", marshalErr)
			d.Nack(false, true)
			return
		}

		outboxEventID := uuid.New().String()
		saveErr := w.repo.SaveOutboxEventTx(ctx, tx, outboxEventID, events.SeatReservationCanceledEventType, payload)
		if saveErr != nil {
			logger.Printf(ctx, "[Meetups Consumer] Failed to save compensation outbox event: %v", saveErr)
			d.Nack(false, true)
			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			logger.Printf(ctx, "[Meetups Consumer] Failed to commit tx for compensation: %v", commitErr)
			d.Nack(false, true)
			return
		}
		d.Ack(false)
		return
	}

	logger.Printf(ctx, "[Meetups Consumer] Internal error during cancel for Workflow %s: %v", cmd.WorkflowID, err)
	d.Nack(false, true)
}
