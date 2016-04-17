package protocol

import (
	"io"
	"net/url"
	"os"
)

type File struct {
}

func init() {
	MustRegister("file", &File{})
}

func (f *File) Open(u *url.URL) (io.ReadCloser, error) {
	// TODO: consider relative path.
	return os.Open(u.Path)
}
