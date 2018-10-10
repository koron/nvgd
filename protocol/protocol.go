package protocol

import (
	"fmt"
	"io"
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
