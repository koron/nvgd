package filterbase_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/filterbase"
)

func TestReadLine(t *testing.T) {
	buf := make([]byte, 16384)
	for i := range buf {
		buf[i] = 'a'
	}
	buf[len(buf)-1] = '\n'

	b := &filterbase.Base{}
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

	if err := b.Close(); err != nil {
		t.Errorf("failed to close: %s", err)
	}
}

func TestConfigMaxLineLen(t *testing.T) {
	if _, err := config.LoadConfig("testdata/maxlen_4K.yml"); err != nil {
		t.Fatal(err)
	}
	if got, want := filterbase.Config.MaxLineLen, 4096; got != want {
		t.Errorf("incorrect max_line_len: want=%d got=%d", want, got)
	}

	if _, err := config.LoadConfig("testdata/maxlen_2M.yml"); err != nil {
		t.Fatal(err)
	}
	if got, want := filterbase.Config.MaxLineLen, 2097152; got != want {
		t.Errorf("incorrect max_line_len: want=%d got=%d", want, got)
	}
}
