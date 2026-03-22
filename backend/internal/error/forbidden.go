package error

import (
	"net/http"
)

type ForbiddenError struct {
	err error
}

func NewForbiddenError(err error) error {
	return &ForbiddenError{
		err: err,
	}
}

func (e *ForbiddenError) Error() string {
	return "Forbidden: user not allowed to perform this action"
}

func (e *ForbiddenError) HTTPStatus() int {
	return http.StatusForbidden
}

func (e *ForbiddenError) Unwrap() error {
	return e.err
}
