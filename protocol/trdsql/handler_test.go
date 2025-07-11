package trdsql

import (
	"io"
	"path"
	"testing"

	"github.com/koron/nvgd/internal/assert"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/internal/protocoltest"
)

func TestRegistered(t *testing.T) {
	protocoltest.GetRegistered[*embedresource.EmbedResource](t, "trdsql")
}

func loadAsset(t *testing.T, name string) string {
	t.Helper()
	f, err := assetFS.Open(name)
	if err != nil {
		t.Fatalf("asset not found: %s", err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		t.Fatalf("failed to load the asset %s: %s", name, err)
	}
	return string(b)
}

func testAsset(t *testing.T, requrl, assetpath string) {
	t.Helper()
	want := loadAsset(t, path.Join("assets", assetpath))
	rsrc := protocoltest.Open(t, requrl)
	got := protocoltest.ReadAllString(t, rsrc)
	assert.Equal(t, want, got, "content mismatch")
}

func TestContents(t *testing.T) {
	testAsset(t, "trdsql:///index.html", "index.html")
	testAsset(t, "trdsql:///editor.css", "editor.css")
	testAsset(t, "trdsql:///editor.js", "editor.js")
}
