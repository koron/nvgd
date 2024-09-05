package filterbase_test

import (
	"testing"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/filterbase"
)

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
