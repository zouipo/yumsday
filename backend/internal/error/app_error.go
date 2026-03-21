package error

import (
	"fmt"
	"log/slog"
	"net/http"
)

// AppError represents the base application error
type AppError struct {
	Message    string
	StatusCode int   // HTTP status code
	Err        error // Underlying error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

/*** COMMON ERROR CONSTRUCTORS ***/

func NewEntityNotFoundError(entityType string, identifier string, err error) error {
	return &AppError{
		Message:    fmt.Sprintf("%s %s not found", entityType, identifier),
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

func NewValidationError(field string, message string, err error) error {
	return &AppError{
		Message:    fmt.Sprintf("Validation error on field '%s': %s", field, message),
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

func NewUnauthorizedError(err error, message string) error {
	return &AppError{
		Message:    fmt.Sprintf("Unauthorized: %s", message),
		StatusCode: http.StatusUnauthorized,
		Err:        err,
	}
}

func NewForbiddenError(err error) error {
	return &AppError{
		Message:    "Forbidden: user not allowed to perform this action",
		StatusCode: http.StatusForbidden,
		Err:        err,
	}
}

func NewConflictError(resource, conflict string, err error) error {
	return &AppError{
		Message:    fmt.Sprintf("Conflict with %s: %s", resource, conflict),
		StatusCode: http.StatusConflict,
		Err:        err,
	}
}

func NewInternalServerError(message string, err error) error {
	if err != nil {
		slog.Error(fmt.Sprintf("%s: %v", message, err))
	} else {
		slog.Error(fmt.Sprintf("%s", message))
	}

	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
