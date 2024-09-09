package echarts

import (
	"io"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/internal/protocoltest"
)

func TestRegistered(t *testing.T) {
	protocoltest.GetRegistered[*embedresource.EmbedResource](t, "echarts")
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
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("content mismatch: -want +got\n%s", d)
	}
}

func TestContents(t *testing.T) {
	testAsset(t, "echarts:///index.html", "index.html")
	testAsset(t, "echarts:///editor.css", "editor.css")
	testAsset(t, "echarts:///editor.js", "editor.js")
}
