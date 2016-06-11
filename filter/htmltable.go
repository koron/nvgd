package filter

import "io"

func newHTMLTable(r io.ReadCloser, p Params) (io.ReadCloser, error) {
	// TODO:
	return r, nil
}

func init() {
	MustRegister("htmltable", newHTMLTable)
}
