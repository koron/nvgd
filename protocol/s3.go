package protocol

import (
	"io"
	"net/url"
)

func init() {
	MustRegister("s3", &S3{})
}

// S3 is AWS S3 protocol handler
type S3 struct {
}

// Open opens a S3 URL.
func (s3 *S3) Open(u *url.URL) (io.ReadCloser, error) {
	// TODO:
	return nil, nil
}
