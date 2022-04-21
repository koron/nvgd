package tail

import (
	"fmt"
	"io"
)

type readSeekCloser interface {
	io.Reader
	io.Seeker
	io.Closer
}

// RTail implements reversed tail filter
type RTail struct {
	raw     readSeekCloser
	limit   int
	bufsize int64
	closed  bool
	ft      bool
}

// NewRTail creates an instance of RTail.
func NewRTail(r readSeekCloser, limit int, bufsize int) *RTail {
	return &RTail{
		raw:     r,
		limit:   limit,
		bufsize: int64(bufsize),
	}
}

// Read reads filtered data.
func (rt *RTail) Read(buf []byte) (int, error) {
	if !rt.ft {
		rt.ft = true
		err := rt.findTop()
		if err != nil {
			return 0, nil
		}
	}
	return rt.raw.Read(buf)
}

// Close closes a filter.
func (rt *RTail) Close() error {
	if rt.closed {
		return nil
	}
	rt.closed = true
	return rt.raw.Close()
}

func (rt *RTail) findTop() error {
	sz, err := rt.raw.Seek(0, 2)
	if err != nil {
		return err
	}
	var (
		cnt  = 0
		buf  = make([]byte, rt.bufsize)
		wbuf []byte
		tail = true
	)
	for curr := sz; curr > 0; {
		if curr < rt.bufsize {
			wbuf = buf[0:curr]
		} else {
			wbuf = buf[0:rt.bufsize]
		}
		curr -= int64(len(wbuf))
		err := rt.seekStart(curr)
		if err != nil {
			return err
		}
		n, err := rt.raw.Read(wbuf)
		if err != nil {
			return err
		}
		if n != len(wbuf) {
			return fmt.Errorf("tail requires %d bytes but got %d", len(wbuf), n)
		}
		for off := n - 1; off >= 0; off-- {
			if wbuf[off] != '\n' || tail {
				tail = false
				continue
			}
			cnt++
			if cnt < rt.limit {
				continue
			}
			return rt.seekStart(curr + int64(off) + 1)
		}
	}
	return rt.seekStart(0)
}

func (rt *RTail) seekStart(off int64) error {
	_, err := rt.raw.Seek(off, 0)
	return err
}
