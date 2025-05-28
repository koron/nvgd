package rangereader

import "io"

type RangeReader struct {
	r io.Reader
	c io.Closer
}

var _ io.ReadCloser = (*RangeReader)(nil)

func New(base io.Reader, start, end, limit int) (*RangeReader, error) {
	r := base
	if start > 0 {
		if s, ok := r.(io.Seeker); ok {
			_, err := s.Seek(int64(start), io.SeekStart)
			if err != nil {
				return nil, err
			}
		} else {
			_, err := io.CopyN(io.Discard, r, int64(start))
			if err != nil {
				return nil, err
			}
		}
	}
	if end < limit-1 {
		r = io.LimitReader(r, int64(end-start+1))
	}
	return &RangeReader{
		r: r,
		c: base.(io.Closer),
	}, nil
}

func (r *RangeReader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *RangeReader) Close() error {
	if r.c == nil {
		return nil
	}
	return r.c.Close()
}
