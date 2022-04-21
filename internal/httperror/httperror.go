package httperror

import "fmt"

type httpErr struct {
	status int
	format string
	args   []interface{}
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

func (err httpErr) statusCode() int {
	return err.status
}

func (err httpErr) body() string {
	return err.Error()
}
