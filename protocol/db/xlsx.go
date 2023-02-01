package db

import (
	"io"

	"github.com/tealeg/xlsx"
)

func openXLSX(r io.Reader) (*xlsx.File, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	xf, err := xlsx.OpenBinary(b)
	if err != nil {
		return nil, err
	}
	return xf, nil
}
