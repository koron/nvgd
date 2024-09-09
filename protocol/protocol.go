// Package protocol provides fundamentals for each protocols.
package protocol

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/koron/nvgd/resource"
)

// Protocol is abstraction of methods to get source stream.
type Protocol interface {
	Open(u *url.URL) (*resource.Resource, error)
}

// ProtocolFunc is Protocol wrapper for function.
type ProtocolFunc func(*url.URL) (*resource.Resource, error)

// Open opens URL as protocol.
func (f ProtocolFunc) Open(u *url.URL) (*resource.Resource, error) {
	return f(u)
}

// Postable is set of methods for POST acceptable source/protocol.
type Postable interface {
	Protocol
	Post(u *url.URL, r io.Reader) (*resource.Resource, error)
}

var protocols = map[string]Protocol{}

// Register registers a protocol with name.
func Register(name string, p Protocol) error {
	_, ok := protocols[name]
	if ok {
		return fmt.Errorf("duplicated protocol name %q", name)
	}
	protocols[name] = p
	return nil
}

// MustRegister registers a protocol, panic if failed.
func MustRegister(name string, p Protocol) {
	if err := Register(name, p); err != nil {
		panic(err)
	}
}

// Find finds a protocol.
func Find(name string) Protocol {
	p, ok := protocols[name]
	if !ok {
		return nil
	}
	return p
}

func Open(u *url.URL, req *http.Request) (*resource.Resource, error) {
	p := Find(u.Scheme)
	if p == nil {
		return nil, fmt.Errorf("not found protocol for %q", u.Scheme)
	}
	if post, ok := p.(Postable); ok && req != nil && req.Method == http.MethodPost {
		return openPost(post, u, req)
	}
	return p.Open(u)
}

func openPost(p Postable, u *url.URL, req *http.Request) (*resource.Resource, error) {
	defer req.Body.Close()
	data := req.Body
	// If the body is multi-part, only the "file00" file is extracted and used.
	// Otherwise, if it is a single stream, it is used as is.
	err := req.ParseMultipartForm(32 * 1024 * 1024) // 32MB
	if err == nil {
		fh, ok := req.MultipartForm.File["file00"]
		if !ok || len(fh) < 1 {
			return nil, errors.New("no files uploaded")
		}
		f, err := fh[0].Open()
		if err != nil {
			return nil, err
		}
		defer f.Close()
		data = f
	} else if !errors.Is(err, http.ErrNotMultipart) {
		return nil, err
	}
	return p.Post(u, data)
}
