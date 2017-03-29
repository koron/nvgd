package db

import (
	"errors"
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
	return nil, errors.New("RestoreHandler: not implemented yet")
}
