package db

import (
	"errors"
	"net/url"

	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

type DumpHandler struct {
}

func init() {
	protocol.MustRegister("db-dump", &DumpHandler{})
}

func (dh *DumpHandler) Open(u *url.URL) (*resource.Resource, error) {
	return nil, errors.New("DumpHandler: not implemented yet")
}
