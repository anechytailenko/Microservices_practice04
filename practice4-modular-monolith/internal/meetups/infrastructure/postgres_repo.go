package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
)

// implement interface Repository defined in the application layer
type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, meetup *domain.Meetup) error {
	query := `
		INSERT INTO meetups (id, title, capacity, status, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		meetup.ID,
		meetup.Title,
		meetup.Capacity,
		meetup.Status,
		meetup.CreatedAt,
	)

	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id domain.MeetupID) (*domain.Meetup, error) {
	query := `
		SELECT id, title, capacity, status, created_at 
		FROM meetups 
		WHERE id = $1`

	var m domain.Meetup
	var rawStatus string

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&m.ID,
		&m.Title,
		&m.Capacity,
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
		SET title = $1, capacity = $2, status = $3 
		WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query,
		meetup.Title,
		meetup.Capacity,
		meetup.Status,
		meetup.ID,
	)

	return err
}
