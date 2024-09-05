package filterbase

import (
	"bufio"
	"errors"
	"io"

	"github.com/koron/nvgd/internal/ltsv"
)

type LTSVReader struct {
	r ltsv.Reader
}

func NewLTSVReader(r io.Reader) *ltsv.Reader {
	return ltsv.NewReaderSize(r, Config.MaxLineLen)
}

func (r *LTSVReader) Read() (*ltsv.Set, error) {
	set, err := r.r.Read()
	if errors.Is(err, bufio.ErrBufferFull) {
		return nil, ErrMaxLineExceeded
	}
	return set, nil
}
