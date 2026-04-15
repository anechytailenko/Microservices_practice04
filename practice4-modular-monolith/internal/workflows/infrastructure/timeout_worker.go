package infrastructure

import (
	"context"
	"database/sql"
	"log"
	"time"
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

	log.Printf("[Workflow Timeout Worker] Started. Checking for sagas older than %v every %v...", w.timeout, pollInterval)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Workflow Timeout Worker] Shutting down...")
			return
		case <-ticker.C:
			w.processTimeouts(ctx)
		}
	}
}

func (w *TimeoutWorker) processTimeouts(ctx context.Context) {
	stuckIDs, err := w.repo.GetStuckWorkflows(ctx, w.timeout)
	if err != nil {
		log.Printf("[Workflow Timeout Worker] Error fetching stuck workflows: %v", err)
		return
	}

	if len(stuckIDs) == 0 {
		return
	}

	for _, id := range stuckIDs {
		err := w.repo.MarkAsTimedOut(ctx, id)
		if err != nil {
			log.Printf("[Workflow Timeout Worker] Failed to mark workflow %s as timed out: %v", id, err)
			continue
		}

		log.Printf("[CRITICAL] Saga %s Time out ->  Moved to ManualIntervention.", id)
	}
}
