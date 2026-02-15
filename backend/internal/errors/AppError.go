package errors

import (
	"fmt"
	"net/http"
)

// AppError represents the base application error
type AppError struct {
	Message    string
	StatusCode int   // HTTP status code
	Err        error // Underlying error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

/*** COMMON ERRORS CONSTRUCTORS ***/

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

func NewUnauthorizedError(err error) error {
	return &AppError{
		Message:    "Unauthorized: authentication required",
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
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
