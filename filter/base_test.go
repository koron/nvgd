package filter

import (
	"bytes"
	"io"
	"testing"
)

func TestBaseReadLine(t *testing.T) {
	buf := make([]byte, 16384)
	for i := range buf {
		buf[i] = 'a'
	}
	buf[len(buf)-1] = '\n'

	b := &Base{}
	b.Init(io.NopCloser(bytes.NewReader(buf)), nil)

	d, err := b.ReadLine()
	if err != nil {
		t.Fatalf("ReadLine failed: %s", err)
	}
	if len(d) != 16384 {
		t.Fatalf("unexpected length: %d", len(d))
	}
	if len(d) != len(buf) {
		t.Fatalf("length not match: %d != %d", len(d), len(buf))
	}
	for i := range d {
		if d[i] != buf[i] {
			t.Fatalf("at %d not match: %c != %c", i, d[i], buf[i])
		}
	}
}
