package protocol

import (
	"io"
	"net/url"
	"os"
)

// File is file protocol handler.
type File struct {
}

func init() {
	MustRegister("file", &File{})
}

// Open opens a URL as file.
func (f *File) Open(u *url.URL) (io.ReadCloser, error) {
	// TODO: consider relative path.
	return os.Open(u.Path)
}
