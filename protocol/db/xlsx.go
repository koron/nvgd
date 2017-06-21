package db

import (
	"io"
	"io/ioutil"

	"github.com/tealeg/xlsx"
)

func openXLSX(r io.Reader) (*xlsx.File, error) {
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
