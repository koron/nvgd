package filterbase

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/koron/nvgd/config"
)

var (
	ErrMaxLineExceeded = errors.New("maximum line length is exceeded. this limit can be extended with config.filters._base_.max_line_len")
)

// Base is base of filters.  It provides common features for filter.
type Base struct {
	Reader *bufio.Reader
	buf    bytes.Buffer
	closed bool
	raw    io.ReadCloser
	rn     BaseReadNext
}

// BaseReadNext is callback to read next data hunk to buf
type BaseReadNext func(buf *bytes.Buffer) error

var Config = struct {
	MaxLineLen int `yaml:"max_line_len"`
}{
	MaxLineLen: 1 * 1024 * 1024,
}

func init() {
	config.RegisterFilter("_base_", &Config)
}

// Init initializes Base object.
func (b *Base) Init(r io.ReadCloser, readNext BaseReadNext) {
	b.Reader = bufio.NewReaderSize(r, Config.MaxLineLen)
	b.raw = r
	b.rn = readNext
}

func (b *Base) Read(buf []byte) (int, error) {
	if b.buf.Len() == 0 {
		if b.closed {
			return 0, io.EOF
		}
		// read next data to b.buf
		b.buf.Reset()
		err := b.rn(&b.buf)
		if err == io.EOF {
			b.Close()
		} else if err != nil {
			return 0, err
		}
	}
	return b.buf.Read(buf)
}

// Close closes head filter.
func (b *Base) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true
	return b.raw.Close()
}
