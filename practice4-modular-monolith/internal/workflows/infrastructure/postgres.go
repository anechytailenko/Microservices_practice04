package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
	"github.com/google/uuid"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, workflow *domain.Workflow, eventID string, eventType string, eventPayload []byte) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	queryWorkflow := `
        INSERT INTO workflows (workflow_id, type, state, last_error, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, queryWorkflow,
		string(workflow.ID),
		workflow.Type,
		string(workflow.State),
		workflow.LastError,
		workflow.CreatedAt,
		workflow.UpdatedAt,
	)
	if err != nil {
		return err
	}

	queryOutbox := `
        INSERT INTO outbox_events (id, event_type, payload, created_at)
        VALUES ($1, $2, $3, $4)`

	_, err = tx.ExecContext(ctx, queryOutbox,
		eventID,
		eventType,
		eventPayload,
		time.Now().UTC(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepo) GetByIDTx(ctx context.Context, tx *sql.Tx, id domain.WorkflowID) (*domain.Workflow, error) {
	query := `
        SELECT workflow_id, type, state, last_error, created_at, updated_at 
        FROM workflows 
        WHERE workflow_id = $1 
        FOR UPDATE`

	var w domain.Workflow
	var rawState string
	var nullLastError sql.NullString

	err := tx.QueryRowContext(ctx, query, string(id)).Scan(
		&w.ID,
		&w.Type,
		&rawState,
		&nullLastError,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.NewNotFoundError("workflow %s not found", string(id))
	}
	if err != nil {
		return nil, err
	}

	w.State = domain.WorkflowState(rawState)
	if nullLastError.Valid {
		w.LastError = &nullLastError.String
	}

	return &w, nil
}

func (r *PostgresRepo) UpdateTx(ctx context.Context, tx *sql.Tx, w *domain.Workflow) error {
	query := `
        UPDATE workflows 
        SET state = $1, last_error = $2, updated_at = $3 
        WHERE workflow_id = $4`

	_, err := tx.ExecContext(ctx, query, string(w.State), w.LastError, w.UpdatedAt, string(w.ID))
	return err
}

func (r *PostgresRepo) SaveOutboxEventTx(ctx context.Context, tx *sql.Tx, eventType string, payload []byte) error {
	query := `INSERT INTO outbox_events (id, event_type, payload, processed, created_at) 
              VALUES ($1, $2, $3, false, NOW())`
	_, err := tx.ExecContext(ctx, query, uuid.New().String(), eventType, payload)
	return err
}

func (r *PostgresRepo) GetStuckWorkflows(ctx context.Context, timeout time.Duration) ([]string, error) {
	query := `
        SELECT workflow_id 
        FROM workflows 
        WHERE state NOT IN ('Completed', 'Failed', 'ManualIntervention') 
        AND updated_at < NOW() - $1::interval`

	intervalStr := fmt.Sprintf("%f seconds", timeout.Seconds())

	rows, err := r.db.QueryContext(ctx, query, intervalStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *PostgresRepo) MarkAsTimedOut(ctx context.Context, workflowID string) error {
	query := `
        UPDATE workflows 
        SET state = 'ManualIntervention', last_error = 'Saga timed out', updated_at = NOW() 
        WHERE workflow_id = $1 
        AND state NOT IN ('Completed', 'Failed', 'ManualIntervention')`

	_, err := r.db.ExecContext(ctx, query, workflowID)
	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id domain.WorkflowID) (*domain.Workflow, error) {
	query := `
        SELECT workflow_id, type, state, last_error, created_at, updated_at 
        FROM workflows 
        WHERE workflow_id = $1`

	var w domain.Workflow
	var rawState string
	var nullLastError sql.NullString

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&w.ID,
		&w.Type,
		&rawState,
		&nullLastError,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	w.State = domain.WorkflowState(rawState)
	if nullLastError.Valid {
		w.LastError = &nullLastError.String
	}
	return &w, nil
}
