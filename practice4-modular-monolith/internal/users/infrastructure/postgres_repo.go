package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/anechytailenko/Microservices_practice04/internal/users/domain"
)

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (id, first_name, last_name, email)
		VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query,
		string(user.ID),
		user.FirstName,
		user.LastName,
		user.Email,
	)

	return err
}

func (r *PostgresRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	query := `
		SELECT id, first_name, last_name, email
		FROM users
		WHERE id = $1`

	var u domain.User

	err := r.db.QueryRowContext(ctx, query, string(id)).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}
