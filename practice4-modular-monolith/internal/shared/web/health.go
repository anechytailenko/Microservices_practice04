package web

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/logger"
)

type HealthResponse struct {
	Status  string            `json:"status"`
	Service string            `json:"service"`
	Commit  string            `json:"commit"`
	Reason  string            `json:"reason,omitempty"`
	Checks  map[string]string `json:"checks,omitempty"`
}

func RegisterHealthRoutes(mux *http.ServeMux, db *sql.DB, serviceName string, commitHash string) {

	mux.HandleFunc("GET /health/live", func(w http.ResponseWriter, r *http.Request) {
		resp := HealthResponse{
			Status:  "UP",
			Service: serviceName,
			Commit:  commitHash,
		}

		WriteJSON(w, http.StatusOK, resp)
	})

	mux.HandleFunc("GET /health/ready", func(w http.ResponseWriter, r *http.Request) {

		if db == nil {
			resp := HealthResponse{
				Status:  "UP",
				Service: serviceName,
				Commit:  commitHash,
				Checks: map[string]string{
					"database": "skipped (no database attached)",
				},
			}
			WriteJSON(w, http.StatusOK, resp)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {

			logger.Printf(ctx, "[%s] Health check failed: db is down: %v", serviceName, err)

			resp := HealthResponse{
				Status:  "DOWN",
				Service: serviceName,
				Commit:  commitHash,
				Reason:  "database unreachable",
			}

			WriteJSON(w, http.StatusServiceUnavailable, resp)
			return
		}

		resp := HealthResponse{
			Status:  "UP",
			Service: serviceName,
			Commit:  commitHash,
			Checks: map[string]string{
				"database": "OK",
			},
		}

		WriteJSON(w, http.StatusOK, resp)
	})
}
