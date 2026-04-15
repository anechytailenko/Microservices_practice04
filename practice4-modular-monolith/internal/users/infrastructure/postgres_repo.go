package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/users/domain"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, user *domain.User) error {
	query := `
        INSERT INTO users (id, first_name, last_name, email, meetups)
        VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		string(user.ID),
		user.FirstName,
		user.LastName,
		user.Email,
		pq.Array(user.Meetups),
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return shared.NewValidationError("user with this email already exists")
		}
		return err
	}

	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query := `
        SELECT id, first_name, last_name, email, meetups
        FROM users
        WHERE id = $1`

	var u domain.User

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		pq.Array(&u.Meetups),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresRepo) GetByIDTx(ctx context.Context, tx *sql.Tx, id string) (*domain.User, error) {
	query := `
        SELECT id, first_name, last_name, email, meetups
        FROM users
        WHERE id = $1 
		FOR UPDATE`

	var u domain.User

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		pq.Array(&u.Meetups),
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, shared.NewNotFoundError("user with id %s not found", id)
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresRepo) UpdateUserMeetupsTx(ctx context.Context, tx *sql.Tx, user *domain.User) error {
	query := `
		UPDATE users 
		SET meetups = $1 
		WHERE id = $2`

	_, err := tx.ExecContext(ctx, query, pq.Array(user.Meetups), string(user.ID))
	return err
}

func (r *PostgresRepo) SaveOutboxEventTx(ctx context.Context, tx *sql.Tx, eventType string, payload []byte) error {
	query := `
		INSERT INTO outbox_events (id, event_type, payload, processed, created_at)
		VALUES ($1, $2, $3, false, NOW())`

	_, err := tx.ExecContext(ctx, query, uuid.New().String(), eventType, payload)
	return err
}
