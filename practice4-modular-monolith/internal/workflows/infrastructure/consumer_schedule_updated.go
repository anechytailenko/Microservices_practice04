package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) handleScheduleUpdated(ctx context.Context, d amqp.Delivery) {
	var evt events.ScheduleUpdated
	if err := json.Unmarshal(d.Body, &evt); err != nil {
		logger.Printf(ctx, "[Workflow Consumer] Failed to unmarshal ScheduleUpdated event: %v", err)
		d.Nack(false, false)
		return
	}

	tx, err := w.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Printf(ctx, "[Workflow Consumer] Failed to begin tx for workflow %s: %v", evt.WorkflowID, err)
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
		logger.Printf(ctx, "[Workflow Consumer] Failed to get workflow %s: %v", evt.WorkflowID, err)
		d.Nack(false, true)
		return
	}

	err = wf.ChangeState(domain.StateCompleted)
	if err != nil {
		var sharedErr shared.Error
		if errors.As(err, &sharedErr) && sharedErr.Type == shared.ErrorTypeConflict {
			logger.Printf(ctx, "[Workflow Consumer] Dropping late event for %s: %v", evt.WorkflowID, err)
			d.Ack(false)
			return
		}
		logger.Printf(ctx, "[Workflow Consumer] Failed to change state for workflow %s: %v", evt.WorkflowID, err)
		d.Nack(false, true)
		return
	}

	err = w.repo.UpdateTx(ctx, tx, wf)
	if err != nil {
		logger.Printf(ctx, "[Workflow Consumer] Failed to update workflow %s: %v", evt.WorkflowID, err)
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
		logger.Printf(ctx, "[Workflow Consumer] Failed to save outbox event for workflow %s: %v", evt.WorkflowID, err)
		d.Nack(false, true)
		return
	}

	if err := tx.Commit(); err != nil {
		logger.Printf(ctx, "[Workflow Consumer] Failed to commit tx for workflow %s: %v", evt.WorkflowID, err)
		d.Nack(false, true)
		return
	}

	logger.Printf(ctx, "[Workflow] Saga %s: SUCCESSFULLY COMPLETED", evt.WorkflowID)
	d.Ack(false)
}
