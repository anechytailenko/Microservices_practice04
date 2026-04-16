package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/notifications/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
	"github.com/jmoiron/sqlx"
)

type notificationDAO struct {
	EventID       string    `db:"event_id"`
	CorrelationID string    `db:"correlation_id"`
	OwnerUserID   string    `db:"owner_user_id"`
	Payload       []byte    `db:"payload"`
	ReceivedAt    time.Time `db:"received_at"`
}

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) SaveNotification(ctx context.Context, eventID, correlationID, ownerUserID string, payload []byte) error {
	query := `
		INSERT INTO notifications (event_id, correlation_id, owner_user_id, payload, received_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (event_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, eventID, correlationID, ownerUserID, payload, time.Now().UTC())
	return err
}

func (r *PostgresRepo) GetByOwnerID(ctx context.Context, ownerUserID string) ([]domain.Notification, error) {
	query := `
		SELECT event_id, correlation_id, owner_user_id, payload, received_at
		FROM notifications
		WHERE owner_user_id = $1
		ORDER BY received_at DESC`

	var daos []notificationDAO
	if err := r.db.SelectContext(ctx, &daos, query, ownerUserID); err != nil {
		return nil, err
	}

	var domainNotifications []domain.Notification

	for _, dao := range daos {
		var rawEvent struct {
			CoreItemId string `json:"coreItemId"`
			Title      string `json:"title"`
			Summary    string `json:"summary"`
			Capacity   int    `json:"capacity"`
			Status     string `json:"status"`
		}

		if err := json.Unmarshal(dao.Payload, &rawEvent); err != nil {
			logger.Println(ctx, "[Repo] Failed to unmarshal payload for event %s: %v", dao.EventID, err)
			continue
		}

		domainNotifications = append(domainNotifications, domain.Notification{
			ID:          dao.EventID,
			OwnerUserID: dao.OwnerUserID,
			MeetupID:    rawEvent.CoreItemId,
			Title:       rawEvent.Title,
			Summary:     rawEvent.Summary,
			Capacity:    rawEvent.Capacity,
			Status:      rawEvent.Status,
			CreatedAt:   dao.ReceivedAt,
		})
	}

	return domainNotifications, nil
}

func (r *PostgresRepo) SaveNotificationTx(ctx context.Context, tx *sql.Tx, eventID, correlationID, userID string, payload []byte) error {
	query := `
		INSERT INTO notifications (event_id, correlation_id, owner_user_id, payload)
		VALUES ($1, $2, $3, $4)`

	_, err := tx.ExecContext(ctx, query, eventID, correlationID, userID, payload)
	return err
}
