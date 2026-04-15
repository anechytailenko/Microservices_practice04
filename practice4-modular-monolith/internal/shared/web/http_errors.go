package web

import (
	"errors"
	"net/http"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func HandleError(w http.ResponseWriter, err error) {
	var sharedErr shared.Error

	status := http.StatusInternalServerError
	message := "Internal Server Error"

	if errors.As(err, &sharedErr) {
		message = sharedErr.Message

		switch sharedErr.Type {
		case shared.ErrorTypeValidation:
			status = http.StatusBadRequest
		case shared.ErrorTypeConflict:
			status = http.StatusConflict
		case shared.ErrorTypeNotFound:
			status = http.StatusNotFound
		case shared.ErrorTypeServiceUnavailable:
			status = http.StatusServiceUnavailable
		case shared.ErrorTypeInternal:
			status = http.StatusInternalServerError
		default:
			status = http.StatusBadRequest
		}
	}

	WriteJSON(w, status, ErrorResponse{Error: message})
}
