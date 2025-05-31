// Package httperror provides error type which composed from HTTP status code
// and message.
package httperror

import (
	"fmt"
	"net/http"
)

type httpErr struct {
	status int
	format string
	args   []interface{}
}

func New(status int) error {
	return Newf(status, http.StatusText(status))
}

func Newf(status int, format string, args ...interface{}) error {
	return &httpErr{
		status: status,
		format: format,
		args:   args,
	}
}

func (err httpErr) Error() string {
	return fmt.Sprintf(err.format, err.args...)
}

func (err httpErr) StatusCode() int {
	return err.status
}

func (err httpErr) Body() string {
	return err.Error()
}
