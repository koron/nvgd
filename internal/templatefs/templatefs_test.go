package templatefs

import (
	"bytes"
	"html/template"
	"io/fs"
	"strings"
	"sync"
	"testing"
	"testing/fstest"
	"time"

	"github.com/koron/nvgd/internal/assert"
)

func TestTemplate(t *testing.T) {
	fsys := fstest.MapFS{
		"hello.tmpl": &fstest.MapFile{
			Data: []byte("Hello, {{.}}!"),
		},
	}
	tfs := New(fsys)
	tmpl, err := tfs.Template("hello.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, "World"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello, World!", buf.String(), "template output")
}

func TestTemplateNotFound(t *testing.T) {
	tfs := New(fstest.MapFS{})
	_, err := tfs.Template("nonexistent.tmpl")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestTemplateIsDir(t *testing.T) {
	fsys := fstest.MapFS{
		"dir": &fstest.MapFile{
			Mode: fs.ModeDir,
		},
	}
	tfs := New(fsys)
	_, err := tfs.Template("dir")
	if err == nil {
		t.Fatal("expected error for directory")
	}
}

func TestTemplateParseError(t *testing.T) {
	fsys := fstest.MapFS{
		"bad.tmpl": &fstest.MapFile{
			Data: []byte("Hello, {{.Name!"),
		},
	}
	tfs := New(fsys)
	_, err := tfs.Template("bad.tmpl")
	if err == nil {
		t.Fatal("expected parse error for invalid template")
	}
}

func TestTemplateWithOption(t *testing.T) {
	fsys := fstest.MapFS{
		"greet.tmpl": &fstest.MapFile{
			Data: []byte("Hello, {{upper .}}!"),
		},
	}
	tfs := New(fsys)
	tmpl, err := tfs.Template("greet.tmpl", OptionFunc(func(tmpl *template.Template) (*template.Template, error) {
		return tmpl.Funcs(map[string]any{
			"upper": strings.ToUpper,
		}), nil
	}))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, "world"); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Hello, WORLD!", buf.String(), "template with FuncMap option")
}

func TestTemplateCache(t *testing.T) {
	now := time.Now()
	fsys := fstest.MapFS{
		"page.tmpl": &fstest.MapFile{
			Data:    []byte("{{.Count}}"),
			ModTime: now,
		},
	}
	tfs := New(fsys)

	tmpl1, err := tfs.Template("page.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	tmpl2, err := tfs.Template("page.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	var buf1, buf2 bytes.Buffer
	if err := tmpl1.Execute(&buf1, struct{ Count int }{42}); err != nil {
		t.Fatal(err)
	}
	if err := tmpl2.Execute(&buf2, struct{ Count int }{42}); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, buf1.String(), buf2.String(), "cached template should produce same output")
}

func TestTemplateConcurrent(t *testing.T) {
	fsys := fstest.MapFS{
		"concurrent.tmpl": &fstest.MapFile{
			Data: []byte("{{.Msg}}"),
		},
	}
	tfs := New(fsys)

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tmpl, err := tfs.Template("concurrent.tmpl")
			if err != nil {
				t.Error(err)
				return
			}
			var buf bytes.Buffer
			if err := tmpl.Execute(&buf, struct{ Msg string }{"ok"}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
}
