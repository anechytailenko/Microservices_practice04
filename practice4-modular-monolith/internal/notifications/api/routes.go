package api

import (
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/notifications/features/get_notifications"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/web"
)

func RegisterRoutes(
	mux *http.ServeMux,
	getNotificationsHandler *get_notifications.Handler,
) {
	// GET /notifications/{ownerId}
	mux.HandleFunc("GET /{ownerId}", func(w http.ResponseWriter, r *http.Request) {
		q := get_notifications.Query{
			OwnerUserID: r.PathValue("ownerId"),
		}

		dtos, err := getNotificationsHandler.Handle(r.Context(), q)
		if err != nil {
			web.HandleError(w, err)
			return
		}

		web.WriteJSON(w, http.StatusOK, dtos)
	})
}
