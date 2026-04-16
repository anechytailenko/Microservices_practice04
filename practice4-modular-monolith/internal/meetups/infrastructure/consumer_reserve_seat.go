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

func (w *ConsumerWorker) processReserveSeat(ctx context.Context, d amqp.Delivery) {
	var cmd commands.ReserveSeat

	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Invalid JSON for ReserveSeat: %v", err)
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

	inboxEventID := fmt.Sprintf("%s-%s", cmd.WorkflowID, commands.ReserveSeatType)
	if !w.checkInbox(ctx, tx, inboxEventID) {
		logger.Printf(ctx, "[Meetups Consumer] Reserve event %s already processed. Skipping.", inboxEventID)
		d.Ack(false)
		return
	}

	meetup, err := w.repo.GetByIDTx(ctx, tx, cmd.MeetupID)
	if err != nil {
		w.handleReserveErrorFlow(ctx, err, tx, d, cmd)
		return
	}

	err = meetup.AddGuest(cmd.UserID)
	if err != nil {
		w.handleReserveErrorFlow(ctx, err, tx, d, cmd)
		return
	}

	err = w.repo.UpdateGuestsTx(ctx, tx, meetup)
	if err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Failed to update DB: %v", err)
		d.Nack(false, true)
		return
	}

	replyEvent := events.SeatReserved{
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
	err = w.repo.SaveOutboxEventTx(ctx, tx, outboxEventID, events.SeatReservedEventType, payload)
	if err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Failed to save outbox event: %v", err)
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Printf(ctx, "[Meetups Consumer] Commit failed: %v", err)
		d.Nack(false, true)
		return
	}

	logger.Printf(ctx, "[Meetups Consumer] Successfully reserved seat for Workflow %s", cmd.WorkflowID)
	d.Ack(false)
}

func (w *ConsumerWorker) handleReserveErrorFlow(ctx context.Context, err error, tx *sql.Tx, d amqp.Delivery, cmd commands.ReserveSeat) {
	var sharedErr shared.Error

	if errors.As(err, &sharedErr) && (sharedErr.Type == shared.ErrorTypeConflict || sharedErr.Type == shared.ErrorTypeNotFound) {
		logger.Printf(ctx, "[Meetups Consumer] Business rule failed for Workflow %s: %s", cmd.WorkflowID, err.Error())

		replyEvent := events.SeatReservationFailed{
			WorkflowID:    cmd.WorkflowID,
			UserID:        cmd.UserID,
			MeetupID:      cmd.MeetupID,
			Reason:        err.Error(),
			CorrelationID: cmd.CorrelationID,
		}

		payload, marshalErr := json.Marshal(replyEvent)
		if marshalErr != nil {
			logger.Printf(ctx, "[Meetups Consumer] Failed to marshal failure reply event: %v", marshalErr)
			d.Nack(false, true)
			return
		}

		outboxEventID := uuid.New().String()
		saveErr := w.repo.SaveOutboxEventTx(ctx, tx, outboxEventID, events.SeatReservationFailedEventType, payload)
		if saveErr != nil {
			logger.Printf(ctx, "[Meetups Consumer] Failed to save failure outbox event: %v", saveErr)
			d.Nack(false, true)
			return
		}

		if commitErr := tx.Commit(); commitErr != nil {
			logger.Printf(ctx, "[Meetups Consumer] Failed to commit transaction for failure event: %v", commitErr)
			d.Nack(false, true)
			return
		}

		d.Ack(false)
		return
	}

	logger.Printf(ctx, "[Meetups Consumer] Internal error for Workflow %s: %v", cmd.WorkflowID, err)
	d.Nack(false, true)
}
