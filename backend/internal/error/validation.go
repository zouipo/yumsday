package error

import (
	"fmt"
	"net/http"
)

type ValidationError struct {
	Field   string
	Message string
	err     error
}

func NewValidationError(field, message string, err error) error {
	return &ValidationError{
		Field:   field,
		Message: message,
		err:     err,
	}
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s", e.Field, e.Message)
}

func (e *ValidationError) HTTPStatus() int {
	return http.StatusBadRequest
}

func (e *ValidationError) Unwrap() error {
	return e.err
}
