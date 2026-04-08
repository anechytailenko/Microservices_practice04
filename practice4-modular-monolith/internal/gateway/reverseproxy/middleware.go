package reverseproxy

import (
	"net/http"

	"github.com/google/uuid"
)

const CorrelationIDHeader = "X-Correlation-Id"

func CorrelationID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = uuid.New().String()
			r.Header.Set(CorrelationIDHeader, correlationID)
		}

		w.Header().Set(CorrelationIDHeader, correlationID)

		next.ServeHTTP(w, r)
	})
}
