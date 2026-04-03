package shared

import (
	"errors"
	"net/http"
)

func HandleError(w http.ResponseWriter, err error) {

	var sharedErr Error

	if errors.As(err, &sharedErr) {
		switch sharedErr.Type {
		case ErrorTypeValidation:
			http.Error(w, sharedErr.Message, http.StatusBadRequest)
		case ErrorTypeConflict:
			http.Error(w, sharedErr.Message, http.StatusConflict)
		case ErrorTypeNotFound:
			http.Error(w, sharedErr.Message, http.StatusNotFound)
		default:
			http.Error(w, sharedErr.Message, http.StatusBadRequest)
		}
		return
	}

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
