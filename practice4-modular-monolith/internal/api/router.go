package api

import (
	"net/http"

	meetups_api "github.com/anechytailenko/Microservices_practice04/internal/meetups/api"
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

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		shared.WriteJSON(w, http.StatusOK, map[string]string{
			"status": "UP",
		})
	})

	meetups_api.RegisterRoutes(
		mux,
		createHandler,
		changeStatusHandler,
		getMeetupHandler,
	)

	return mux
}
