package tail

import (
	"errors"
	"io"
	"testing"
)

type errReadSeekCloser struct {
	err error
}

func (r *errReadSeekCloser) Read(p []byte) (int, error) {
	return 0, r.err
}

func (r *errReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	return 0, r.err
}

func (r *errReadSeekCloser) Close() error {
	return nil
}

func TestRTailReadError(t *testing.T) {
	want := errors.New("simulated seek error")
	rt := NewRTail(&errReadSeekCloser{err: want}, 10, 4096)
	_, err := rt.Read(nil)
	if err != want {
		t.Fatalf("expected error %v, got %v", want, err)
	}
}

func TestRTailReadErrorNotEOF(t *testing.T) {
	want := errors.New("disk failure")
	rt := NewRTail(&errReadSeekCloser{err: want}, 1, 4096)
	_, err := rt.Read(nil)
	if err == nil || err == io.EOF {
		t.Fatalf("expected non-EOF error, got %v", err)
	}
}
