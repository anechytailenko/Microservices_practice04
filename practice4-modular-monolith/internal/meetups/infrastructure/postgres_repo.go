package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, meetup *domain.Meetup, eventID string, eventType string, eventPayload []byte) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	queryMeetup := `
        INSERT INTO meetups (id, title, capacity, owner_user_id, guests, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err = tx.ExecContext(ctx, queryMeetup,
		string(meetup.ID),
		meetup.Title,
		meetup.Capacity,
		meetup.OwnerUserID,
		pq.Array(meetup.Guests),
		string(meetup.Status),
		meetup.CreatedAt,
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

func (r *PostgresRepo) GetByID(ctx context.Context, id domain.MeetupID) (*domain.Meetup, error) {
	query := `
        SELECT id, title, capacity, owner_user_id, guests, status, created_at 
        FROM meetups 
        WHERE id = $1`

	var m domain.Meetup
	var rawStatus string

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&m.ID,
		&m.Title,
		&m.Capacity,
		&m.OwnerUserID,
		pq.Array(&m.Guests),
		&rawStatus,
		&m.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	m.Status = domain.MeetupStatus(rawStatus)
	return &m, nil
}

func (r *PostgresRepo) Update(ctx context.Context, meetup *domain.Meetup) error {
	query := `
        UPDATE meetups 
        SET title = $1, capacity = $2, owner_user_id = $3, guests = $4, status = $5 
        WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query,
		meetup.Title,
		meetup.Capacity,
		meetup.OwnerUserID,
		pq.Array(meetup.Guests),
		string(meetup.Status),
		string(meetup.ID),
	)

	return err
}

func (r *PostgresRepo) GetByIDTx(ctx context.Context, tx *sql.Tx, id string) (*domain.Meetup, error) {
	query := `
		SELECT id, title, capacity, owner_user_id, guests, status, created_at
		FROM meetups 
		WHERE id = $1 
		FOR UPDATE`

	var m domain.Meetup

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.Title,
		&m.Capacity,
		&m.OwnerUserID,
		pq.Array(&m.Guests),
		&m.Status,
		&m.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.NewNotFoundError("meetup with id '%s' not found", id)
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *PostgresRepo) UpdateGuestsTx(ctx context.Context, tx *sql.Tx, meetup *domain.Meetup) error {
	query := `UPDATE meetups SET guests = $2 WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, string(meetup.ID), pq.Array(meetup.Guests))
	return err
}

func (r *PostgresRepo) SaveOutboxEventTx(ctx context.Context, tx *sql.Tx, eventID string, eventType string, payload []byte) error {
	query := `
		INSERT INTO outbox_events (id, event_type, payload, processed, created_at)
		VALUES ($1, $2, $3, false, NOW())`

	_, err := tx.ExecContext(ctx, query, eventID, eventType, payload)
	return err
}
