package core

type httpError interface {
	StatusCode() int
	Body() string
}
