package api

import (
	"encoding/json"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	shared "github.com/anechytailenko/Microservices_practice04/internal/shared/web"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createHandler *create_meetup.Handler,
	changeStatusHandler *change_status.Handler,
	getMeetupHandler *get_meetup.Handler,
) {
	// POST /
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var cmd create_meetup.Command

		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		corrID := r.Header.Get("X-Correlation-Id")

		ctx := ctxutil.WithCorrelationID(r.Context(), corrID)

		id, err := createHandler.Handle(ctx, cmd)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"id": string(id)})
	})

	// PATCH /{id}/status
	mux.HandleFunc("PATCH /{id}/status", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		var reqBody struct {
			Status string `json:"status"`
		}

		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		cmd := change_status.Command{
			MeetupID: id,
			Status:   reqBody.Status,
		}

		if err := changeStatusHandler.Handle(r.Context(), cmd); err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteNoContent(w)
	})

	// GET /{id}
	mux.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		q := get_meetup.Query{
			MeetupID: r.PathValue("id"),
		}

		dto, err := getMeetupHandler.Handle(r.Context(), q)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusOK, dto)
	})
}
