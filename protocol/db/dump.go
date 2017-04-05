package db

import (
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
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	table := path(u)
	rows, err := c.db.Query("SELECT * FROM " + table)
	if err != nil {
		return nil, err
	}
	rc, err := rows2ltsv(rows, -1)
	if err != nil {
		return nil, err
	}
	return resource.New(rc), nil
}
