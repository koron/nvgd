package filterbase

import (
	"bufio"
	"errors"
	"io"
	"iter"

	"github.com/koron/nvgd/internal/ltsv"
)

type LTSVReader struct {
	r *ltsv.Reader
}

func NewLTSVReader(r io.Reader) *LTSVReader {
	return &LTSVReader{
		r: ltsv.NewReaderSize(r, Config.MaxLineLen),
	}
}

func (r *LTSVReader) Read() (*ltsv.Set, error) {
	set, err := r.r.Read()
	if errors.Is(err, bufio.ErrBufferFull) {
		return nil, ErrMaxLineExceeded
	}
	return set, err
}

func (r *LTSVReader) Iter() iter.Seq2[*ltsv.Set, error] {
	return func(yield func(*ltsv.Set, error) bool) {
		for {
			s, err := r.Read()
			if err != nil {
				if !errors.Is(err, io.EOF) {
					yield(nil, err)
					return
				}
				return
			}
			if s.Empty() {
				continue
			}
			if !yield(s, nil) {
				return
			}
		}
	}
}
