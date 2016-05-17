package core

type httpError interface {
	statusCode() int
	body() string
}
