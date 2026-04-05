package api

import (
	"encoding/json"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/change_status"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/create_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/meetups/features/get_meetup"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

func NewRouter(
	createHandler *create_meetup.Handler,
	changeStatusHandler *change_status.Handler,
	getMeetupHandler *get_meetup.Handler,
) *http.ServeMux {

	mux := http.NewServeMux()

	// GET /health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "UP",
		})
	})

	// POST /meetups
	mux.HandleFunc("POST /meetups", func(w http.ResponseWriter, r *http.Request) {
		var cmd create_meetup.Command

		if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}

		id, err := createHandler.Handle(r.Context(), cmd)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusCreated, map[string]string{"id": string(id)})
	})

	// PATCH /meetups/{id}/status
	mux.HandleFunc("PATCH /meetups/{id}/status", func(w http.ResponseWriter, r *http.Request) {
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

	// GET /meetups/{id}
	mux.HandleFunc("GET /meetups/{id}", func(w http.ResponseWriter, r *http.Request) {
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

	return mux
}
