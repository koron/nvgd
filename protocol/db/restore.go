package db

import (
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"net/url"

	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
	"github.com/tealeg/xlsx"
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
	xf, err := rh.openXLSX(r)
	if err != nil {
		return nil, err
	}
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	for _, xs := range xf.Sheets {
		err := rh.restoreSheet(tx, xs)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return nil, errors.New("RestoreHandler#Post: not implemented yet")
}

func (rh *RestoreHandler) openXLSX(r io.Reader) (*xlsx.File, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	xf, err := xlsx.OpenBinary(b)
	if err != nil {
		return nil, err
	}
	return xf, nil
}

func (rh *RestoreHandler) restoreSheet(tx *sql.Tx, xs *xlsx.Sheet) error {
	// TODO:
	return nil
}
