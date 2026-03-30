package error

import (
	"fmt"
	"net/http"
)

type InvalidParamsError struct {
	Field   string
	Message string
	err     error
}

func NewInvalidParamsError(field string, err error) error {
	return &InvalidParamsError{
		Field: field,
		err:   err,
	}
}

func (e *InvalidParamsError) Error() string {
	return fmt.Sprintf("Invalid parameter '%s'", e.Field)
}

func (e *InvalidParamsError) HTTPStatus() int {
	return http.StatusBadRequest
}

func (e *InvalidParamsError) Unwrap() error {
	return e.err
}
