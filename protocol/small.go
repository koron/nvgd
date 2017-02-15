package protocol

import "io"

type Small interface {
	Small()
}

type SmallReadCloser struct {
	io.Reader
}

func (SmallReadCloser) Close() error { return nil }

func (SmallReadCloser) Small() {}
