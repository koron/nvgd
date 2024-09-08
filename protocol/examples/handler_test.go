package examples

import (
	"io"
	"path"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/koron/nvgd/internal/embedresource"
	"github.com/koron/nvgd/internal/protocoltest"
)

func TestRegistered(t *testing.T) {
	protocoltest.GetRegistered[*embedresource.EmbedResource](t, "examples")
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
	testAsset(t, "examples:///line.csv", "line.csv")
	testAsset(t, "examples:///pie.csv", "pie.csv")
	testAsset(t, "examples:///test.csv", "test.csv")
	// Currently index.html is not available.
	//testAsset(t, "examples:///", "index.html")
}
