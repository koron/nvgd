package resource

import "io"

// ReadSeekCloser combines io.Reader, io.Seeker and io.Closer.
type ReadSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}
