package error

import (
	"fmt"
	"net/http"
	"strings"
)

type InvalidParamsError struct {
	Fields []string
	err    error
}

func NewInvalidParamsError(fields []string, err error) error {
	return &InvalidParamsError{
		Fields: fields,
		err:    err,
	}
}

func (e *InvalidParamsError) Error() string {
	if len(e.Fields) == 1 {
		return fmt.Sprintf("Invalid parameter '%s'", e.Fields[0])
	}
	return fmt.Sprintf("Invalid parameters '%s'", strings.Join(e.Fields, "', '"))
}

func (e *InvalidParamsError) AddInvalidField(field string) {
	e.Fields = append(e.Fields, field)
}

func (e *InvalidParamsError) HTTPStatus() int {
	return http.StatusBadRequest
}

func (e *InvalidParamsError) Unwrap() error {
	return e.err
}
