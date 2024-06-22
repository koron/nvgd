package core

import (
	"errors"
	"io/fs"
	"net/http"
)

type httpError interface {
	StatusCode() int
	Body() string
}

func toHTTPError(err error) (msg string, code int) {
	var herr httpError
	if errors.As(err, &herr) {
		return herr.Body(), herr.StatusCode()
	}
	if errors.Is(err, fs.ErrNotExist) {
		return "404 page not found: " + err.Error(), http.StatusNotFound
	}
	if errors.Is(err, fs.ErrPermission) {
		return "403 Forbidden: " + err.Error(), http.StatusForbidden
	}
	// Default:
	return "500 Internal Server Error: " + err.Error(), http.StatusInternalServerError
}
