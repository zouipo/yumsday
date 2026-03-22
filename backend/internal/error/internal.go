package error

import (
	"fmt"
	"log/slog"
	"net/http"
)

type InternalError struct {
	Message string
	err     error
}

func NewInternalError(message string, err error) error {
	if err == nil {
		slog.Error(fmt.Sprintf("%s", message))
	} else {
		slog.Error(fmt.Sprintf("%s: %v", message, err))
	}

	return &InternalError{
		Message: message,
		err:     err,
	}
}

func (e *InternalError) Error() string {
	return e.Message
}

func (e *InternalError) HTTPStatus() int {
	return http.StatusInternalServerError
}

func (e *InternalError) Unwrap() error {
	return e.err
}
