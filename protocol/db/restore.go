package db

import (
	"errors"
	"io"
	"net/url"

	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

type RestoreHandler struct {
}

func init() {
	protocol.MustRegister("db-restore", &RestoreHandler{})
}

func (rh *RestoreHandler) Open(u *url.URL) (*resource.Resource, error) {
	return nil, errors.New("RestoreHandler#Open: not implemented yet")
}

func (rh *RestoreHandler) Post(u *url.URL, r io.Reader) (*resource.Resource, error) {
	return nil, errors.New("RestoreHandler#Post: not implemented yet")
}
