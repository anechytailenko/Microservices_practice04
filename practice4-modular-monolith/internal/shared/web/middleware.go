package web

import (
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
)

func ContextWithCorrelationID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := r.Header.Get(ctxutil.CorrelationIDHeader)
		ctx := ctxutil.WithCorrelationID(r.Context(), corrID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
