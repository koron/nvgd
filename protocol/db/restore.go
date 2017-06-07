package db

import (
	"errors"
	"io"
	"net/url"

	xlsx4db "github.com/koron/go-xlsx4db"
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
	xf, err := openXLSX(r)
	if err != nil {
		return nil, err
	}
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	tables := parseAsTables(u)
	err = xlsx4db.Restore(c.db, xf, true, tables...)
	if err != nil {
		return nil, err
	}
	return resource.NewString("restored successfully"), nil
}
