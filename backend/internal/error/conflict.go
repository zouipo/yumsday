package error

import (
	"fmt"
	"net/http"
)

type ConflictError struct {
	Entity  string
	Message string
	err     error
}

func NewConflictError(entity, message string, err error) error {
	return &ConflictError{
		Entity:  entity,
		Message: message,
		err:     err,
	}
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("Conflict with entity '%s': %s", e.Entity, e.Message)
}

func (e *ConflictError) HTTPStatus() int {
	return http.StatusConflict
}

func (e *ConflictError) Unwrap() error {
	return e.err
}
