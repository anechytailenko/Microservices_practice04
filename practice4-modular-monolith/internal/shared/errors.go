package shared

import "fmt"

type ErrorType string

const (
	ErrorTypeValidation         ErrorType = "validation"
	ErrorTypeConflict           ErrorType = "conflict"
	ErrorTypeNotFound           ErrorType = "not_found"
	ErrorTypeServiceUnavailable ErrorType = "service_unavailable"
)

type Error struct {
	Type    ErrorType
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func NewValidationError(msg string, args ...any) Error {
	return Error{
		Type:    ErrorTypeValidation,
		Message: fmt.Sprintf(msg, args...),
	}
}

func NewConflictError(msg string, args ...any) Error {
	return Error{
		Type:    ErrorTypeConflict,
		Message: fmt.Sprintf(msg, args...),
	}
}

func NewNotFoundError(msg string, args ...any) Error {
	return Error{
		Type:    ErrorTypeNotFound,
		Message: fmt.Sprintf(msg, args...),
	}
}

func NewServiceUnavailableError(msg string, args ...any) Error {
	return Error{
		Type:    ErrorTypeServiceUnavailable,
		Message: fmt.Sprintf(msg, args...),
	}
}
