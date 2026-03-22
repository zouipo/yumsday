package error

import (
	"fmt"
	"net/http"
)

type UnauthorizedError struct {
	Message string
	err     error
}

func NewUnauthorizedError(message string, err error) error {
	return &UnauthorizedError{
		Message: message,
		err:     err,
	}
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %s", e.Message)
}

func (e *UnauthorizedError) HTTPStatus() int {
	return http.StatusUnauthorized
}

func (e *UnauthorizedError) Unwrap() error {
	return e.err
}
