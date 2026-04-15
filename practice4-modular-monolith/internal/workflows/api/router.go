package api

import (
	"encoding/json"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/features/create_workflow"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/features/get_workflow"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createHandler *create_workflow.Handler,
	getHandler *get_workflow.Handler,
) {
	// POST /join-meetup
	mux.HandleFunc("POST /join-meetup", func(w http.ResponseWriter, r *http.Request) {
		var cmd create_workflow.Command

		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		corrID := r.Header.Get("X-Correlation-Id")
		ctx := ctxutil.WithCorrelationID(r.Context(), corrID)

		dto, err := createHandler.Handle(ctx, cmd)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusCreated, dto)
	})

	// GET /{id}
	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		q := get_workflow.Query{
			WorkflowID: r.PathValue("id"),
		}

		dto, err := getHandler.Handle(r.Context(), q)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusOK, dto)
	})
}
