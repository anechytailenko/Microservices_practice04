package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
)

type TimeoutWorker struct {
	db      *sql.DB
	repo    *PostgresRepo
	timeout time.Duration
}

func NewTimeoutWorker(db *sql.DB, repo *PostgresRepo, timeout time.Duration) *TimeoutWorker {
	return &TimeoutWorker{
		db:      db,
		repo:    repo,
		timeout: timeout,
	}
}

func (w *TimeoutWorker) Start(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	logger.Printf(ctx, "[Workflow Timeout Worker] Started. Checking for sagas older than %v every %v...", w.timeout, pollInterval)

	for {
		select {
		case <-ctx.Done():
			logger.Println(ctx, "[Workflow Timeout Worker] Shutting down...")
			return
		case <-ticker.C:
			w.processTimeouts(ctx)
		}
	}
}

func (w *TimeoutWorker) processTimeouts(ctx context.Context) {
	stuckIDs, err := w.repo.GetStuckWorkflows(ctx, w.timeout)
	if err != nil {
		logger.Printf(ctx, "[Workflow Timeout Worker] Error fetching stuck workflows: %v", err)
		return
	}

	if len(stuckIDs) == 0 {
		return
	}

	for _, id := range stuckIDs {
		err := w.repo.MarkAsTimedOut(ctx, id)
		if err != nil {
			logger.Printf(ctx, "[Workflow Timeout Worker] Failed to mark workflow %s as timed out: %v", id, err)
			continue
		}

		logger.Printf(ctx, "[CRITICAL] Saga %s Time out ->  Moved to ManualIntervention.", id)
	}
}
