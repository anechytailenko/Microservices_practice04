package web

import (
	"errors"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

func HandleError(w http.ResponseWriter, err error) {
	var sharedErr shared.Error

	if errors.As(err, &sharedErr) {
		switch sharedErr.Type {
		case shared.ErrorTypeValidation:
			http.Error(w, sharedErr.Message, http.StatusBadRequest)
		case shared.ErrorTypeConflict:
			http.Error(w, sharedErr.Message, http.StatusConflict)
		case shared.ErrorTypeNotFound:
			http.Error(w, sharedErr.Message, http.StatusNotFound)
		case shared.ErrorTypeServiceUnavailable:
			http.Error(w, sharedErr.Message, http.StatusServiceUnavailable)
		case shared.ErrorTypeInternal:
			http.Error(w, sharedErr.Message, http.StatusInternalServerError)
		default:
			http.Error(w, sharedErr.Message, http.StatusBadRequest)
		}
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
