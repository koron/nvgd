package db

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
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
	for _, sh := range xf.Sheets {
		log.Printf("Sheet: %q", sh.Name)
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
