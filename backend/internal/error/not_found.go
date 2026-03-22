package error

import (
	"fmt"
	"net/http"
)

type NotFoundError struct {
	Entity     string
	Identifier string
	err        error
}

func NewNotFoundError(entity, identifier string, err error) error {
	return &NotFoundError{
		Entity:     entity,
		Identifier: identifier,
		err:        err,
	}
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("Entity '%s' with identifier '%s' not found", e.Entity, e.Identifier)
}

func (e *NotFoundError) HTTPStatus() int {
	return http.StatusNotFound
}

func (e *NotFoundError) Unwrap() error {
	return e.err
}
