package help

import (
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/nvgd/doc"
	"github.com/koron/nvgd/internal/protocoltest"
	"github.com/koron/nvgd/protocol"
)

func TestRegistered(t *testing.T) {
	protocoltest.CheckRegistered(t, "help", protocol.ProtocolFunc(Serve))
}

func TestRootText(t *testing.T) {
	testText := func(t *testing.T, want string) {
		t.Helper()
		Text = want
		rsrc := protocoltest.Open(t, "help:///")
		defer rsrc.Close()
		got := protocoltest.ReadAllString(t, rsrc)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("root text unmatch:\nwant=%q\ngot=%q", want, got)
		}
	}
	testText(t, "")
	testText(t, "Hello World")
	testText(t, "") // reset
}

func TestRedirect(t *testing.T) {
	rsrc := protocoltest.Open(t, "help://")
	defer rsrc.Close()
	protocoltest.CheckRedirect(t, rsrc, "help:///")
}

func TestDoc(t *testing.T) {
	readDoc := func(t *testing.T, name string) string {
		f, err := doc.FS.Open(name)
		if err != nil {
			t.Fatalf("failed to open: %s", err)
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			t.Fatalf("failed to read %s: %s", name, err)
		}
		return string(b)
	}
	testExist := func(t *testing.T, requrl, docpath string) {
		t.Helper()
		rsrc := protocoltest.Open(t, requrl)
		defer rsrc.Close()
		protocoltest.CheckNotRedirect(t, rsrc)
		// compare content
		got := protocoltest.ReadAllString(t, rsrc)
		want := readDoc(t, docpath)
		if d := cmp.Diff(want, got); d != "" {
			t.Errorf("contents unmatch %s: +want -got\n%s", docpath, d)
		}
	}
	testNotExist := func(t *testing.T, requrl string) {
		rsrc := protocoltest.Open(t, requrl)
		defer rsrc.Close()
		protocoltest.CheckRedirect(t, rsrc, "help:///")
	}
	testExist(t, "help:///doc/protocol-db.md", "protocol-db.md")
	testExist(t, "help:///doc/protocol-redis.md", "protocol-redis.md")
	testExist(t, "help:///doc/filter-echarts.md", "filter-echarts.md")
	testExist(t, "help:///doc/filter-trdsql.md", "filter-trdsql.md")
	testNotExist(t, "help:///doc/_not_exists_.md")
}
