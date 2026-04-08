package api

import (
	"encoding/json"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/users/features/create_user"
	"github.com/anechytailenko/Microservices_practice04/internal/users/features/get_user"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createHandler *create_user.Handler,
	getUserHandler *get_user.Handler,
) {
	// POST /users
	mux.HandleFunc("POST /users", func(w http.ResponseWriter, r *http.Request) {
		var cmd create_user.Command

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

	// GET /users/{id}
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		q := get_user.Query{
			UserID: r.PathValue("id"),
		}

		dto, err := getUserHandler.Handle(r.Context(), q)
		if err != nil {
			shared.HandleError(w, err)
			return
		}

		shared.WriteJSON(w, http.StatusOK, dto)
	})
}
