package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
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
        INSERT INTO meetups (id, title, capacity, owner_user_id, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, queryMeetup,
		string(meetup.ID),
		meetup.Title,
		meetup.Capacity,
		meetup.OwnerUserID,
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
        SELECT id, title, capacity, owner_user_id, status, created_at 
        FROM meetups 
        WHERE id = $1`

	var m domain.Meetup
	var rawStatus string

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&m.ID,
		&m.Title,
		&m.Capacity,
		&m.OwnerUserID,
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
        SET title = $1, capacity = $2, owner_user_id = $3, status = $4 
        WHERE id = $5`

	_, err := r.db.ExecContext(ctx, query,
		meetup.Title,
		meetup.Capacity,
		meetup.OwnerUserID,
		string(meetup.Status),
		string(meetup.ID),
	)

	return err
}
