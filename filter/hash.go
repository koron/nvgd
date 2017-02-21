package filter

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"strings"

	"github.com/koron/nvgd/resource"
)

// Hash represents hash filter.
type Hash struct {
	Base
	s   int
	h   hash.Hash
	enc hashEncoder
}

type hashEncoder func(io.Writer, []byte) error

// NewHash creates a hash filter instance.
func NewHash(r io.ReadCloser, h hash.Hash, enc hashEncoder) *Hash {
	f := &Hash{
		s:   0,
		h:   h,
		enc: enc,
	}
	f.Base.Init(r, f.readNext)
	return f
}

func (f *Hash) readNext(buf *bytes.Buffer) error {
	if f.s == 2 {
		return io.EOF
	}
	if f.s == 0 {
		if err := f.readAll(); err != nil {
			return err
		}
		f.s = 1
	}
	if err := f.enc(buf, f.h.Sum(nil)); err != nil {
		return err
	}
	f.s = 2
	return nil
}

func (f *Hash) readAll() error {
	_, err := io.Copy(f.h, f.Reader)
	return err
}

func toHash(n string) (hash.Hash, error) {
	switch strings.ToLower(n) {
	case "md5":
		return md5.New(), nil
	case "sha1":
		return sha1.New(), nil
	case "sha256":
		return sha256.New(), nil
	case "sha512":
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("unknown hash algorithm: %s", n)
	}
}

func toEnc(n string) (hashEncoder, error) {
	switch strings.ToLower(n) {
	case "hex":
		return hashEncHex, nil
	case "base64":
		return hashEncBase64, nil
	case "binary":
		return hashEncBin, nil
	default:
		return nil, fmt.Errorf("unknown hash encoder: %s", n)
	}
}

func hashEncHex(w io.Writer, b []byte) error {
	_, err := fmt.Fprintf(w, "%x", b)
	return err
}

func hashEncBase64(w io.Writer, b []byte) error {
	s := base64.StdEncoding.EncodeToString(b)
	_, err := w.Write([]byte(s))
	return err
}

func hashEncBin(w io.Writer, b []byte) error {
	_, err := w.Write(b)
	return err
}

func newHash(r *resource.Resource, p Params) (*resource.Resource, error) {
	h, err := toHash(p.String("algorithm", "md5"))
	if err != nil {
		return nil, err
	}
	enc, err := toEnc(p.String("encoding", "hex"))
	if err != nil {
		return nil, err
	}
	return r.Wrap(NewHash(r, h, enc)), nil
}

func init() {
	MustRegister("hash", newHash)
}
