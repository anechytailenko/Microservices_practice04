package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/events"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *ConsumerWorker) handleSeatCanceled(ctx context.Context, d amqp.Delivery) {
	var evt events.SeatReservationCanceled
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

	inboxID := fmt.Sprintf("%s-%s", evt.WorkflowID, events.SeatReservationCanceledEventType)
	if !w.checkInbox(ctx, tx, inboxID) {
		d.Ack(false)
		return
	}

	wf, err := w.repo.GetByIDTx(ctx, tx, domain.WorkflowID(evt.WorkflowID))
	if err != nil {
		d.Nack(false, true)
		return
	}

	err = wf.ChangeState(domain.StateFailed)
	if err != nil {
		var sharedErr shared.Error
		if errors.As(err, &sharedErr) && sharedErr.Type == shared.ErrorTypeConflict {
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

	if err := tx.Commit(); err != nil {
		d.Nack(false, true)
		return
	}

	log.Printf("[Workflow] Saga %s: Compensation completed. Saga is now FAILED.", evt.WorkflowID)
	d.Ack(false)
}
