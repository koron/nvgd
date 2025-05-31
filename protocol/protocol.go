// Package protocol provides fundamentals for each protocols.
package protocol

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/koron/nvgd/internal/httperror"
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

type Rangeable interface {
	Protocol
	Size(u *url.URL) (int, error)
	OpenRange(u *url.URL, start, end int) (*resource.Resource, error)
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
	if rangeable, ok := p.(Rangeable); ok && req != nil {
		if req.Method == http.MethodHead {
			return openRangeHead(rangeable, u, req)
		}
		if req.Method == http.MethodGet {
			return openRangeBody(rangeable, u, req)
		}
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

func openRangeHead(rangeable Rangeable, u *url.URL, req *http.Request) (*resource.Resource, error) {
	sz, err := rangeable.Size(u)
	if err != nil {
		// TODO: Return better error.
		return nil, fmt.Errorf("failed to fetch size: %w", err)
	}
	// TODO: Prepare better resource.
	r, err := rangeable.Open(u)
	if err != nil {
		return nil, err
	}
	// Add headers to accept range request.
	r.Put(resource.SkipFilters, true)
	r.Put(resource.AcceptRanges, "bytes")
	r.Put(resource.ContentLength, sz)
	return r, nil
}

var rxBytesRange = regexp.MustCompile(`^bytes=(\d+)-(\d+)$`)

func openRangeBody(rangeable Rangeable, u *url.URL, req *http.Request) (*resource.Resource, error) {
	// Parse and extract "Range" header.
	rangeHeader := req.Header.Get("Range")
	if rangeHeader == "" {
		return rangeable.Open(u)
	}
	m := rxBytesRange.FindStringSubmatch(rangeHeader)
	if m == nil {
		return nil, httperror.New(http.StatusBadRequest)
	}
	start, _ := strconv.Atoi(m[1])
	end, _ := strconv.Atoi(m[2])

	// TODO: Check if in range.
	sz, err := rangeable.Size(u)
	if err != nil {
		// TODO: Return better error.
		return nil, fmt.Errorf("failed to fetch size: %w", err)
	}

	r, err := rangeable.OpenRange(u, start, end)
	if err != nil {
		return nil, err
	}
	r.Put(resource.SkipFilters, true)
	r.Put(resource.ContentLength, end-start+1)
	r.Put(resource.ContentRange, fmt.Sprintf("bytes %d-%d/%d", start, end, sz))
	return r, nil
}
