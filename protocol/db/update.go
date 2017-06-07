package db

import (
	"errors"
	"io"
	"net/url"

	xlsx4db "github.com/koron/go-xlsx4db"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

type UpdateHandler struct {
}

func init() {
	protocol.MustRegister("db-update", &UpdateHandler{})
}

func (uh *UpdateHandler) Open(u *url.URL) (*resource.Resource, error) {
	return nil, errors.New("UpdateHandler#Open: not implemented yet")
}

func (uh *UpdateHandler) Post(u *url.URL, r io.Reader) (*resource.Resource, error) {
	xf, err := openXLSX(r)
	if err != nil {
		return nil, err
	}
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	tables := parseAsTables(u)
	err = xlsx4db.Update(c.db, xf, tables...)
	if err != nil {
		return nil, err
	}
	return resource.NewString("updated successfully"), nil
}
