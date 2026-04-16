package reverseproxy

import (
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/google/uuid"
)

func CorrelationID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		correlationID := r.Header.Get(ctxutil.CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
			r.Header.Set(ctxutil.CorrelationIDHeader, correlationID)
		}

		w.Header().Set(ctxutil.CorrelationIDHeader, correlationID)

		next.ServeHTTP(w, r)
	})
}
