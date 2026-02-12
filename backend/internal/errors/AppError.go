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

// NewEntityNotFoundError creates an error indicating that a specific entity was not found.
func NewEntityNotFoundError(entityType string, identifier string, err error) error {
	return &AppError{
		Message:    fmt.Sprintf("%s %s not found", entityType, identifier),
		StatusCode: http.StatusNotFound,
		Err:        err,
	}
}

// NewValidationError creates an error indicating that a validation error occurred on a specific field with a message describing the issue.
func NewValidationError(field string, message string, err error) error {
	return &AppError{
		Message:    fmt.Sprintf("Validation error on field '%s': %s", field, message),
		StatusCode: http.StatusBadRequest,
		Err:        err,
	}
}

// NewUnauthorizedError creates an error indicating that the user is not authenticated.
func NewUnauthorizedError(err error) error {
	return &AppError{
		Message:    "Unauthorized: authentication required",
		StatusCode: http.StatusUnauthorized,
		Err:        err,
	}
}

// NewForbiddenError creates an error indicating that the user is authenticated but does not have permission to perform the action.
func NewForbiddenError(err error) error {
	return &AppError{
		Message:    "Forbidden: user not allowed to perform this action",
		StatusCode: http.StatusForbidden,
		Err:        err,
	}
}

// NewConflictError creates an error indicating that there is a conflict with the current state of the resource, such as a duplicate entry or a version mismatch.
func NewConflictError(resource, conflict string, err error) error {
	return &AppError{
		Message:    fmt.Sprintf("Conflict with %s: %s", resource, conflict),
		StatusCode: http.StatusConflict,
		Err:        err,
	}
}

// NewInternalServerError creates an error indicating that an unexpected internal error occurred, with a message describing the issue.
func NewInternalServerError(message string, err error) error {
	return &AppError{
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
}
