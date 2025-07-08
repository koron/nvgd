// Package toopfs provides UI to download files to the OPFS.
package toopfs

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"io"
	"strconv"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/internal/ltsv"
	"github.com/koron/nvgd/resource"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "filter/toopfs", "")

type doc struct {
	Downloads []entry
}

type entry struct {
	Name string
	Type string
	Size uint64
	Link string
}

func (e entry) IsDir() bool {
	return e.Type == "dir"
}

func toEntry(s *ltsv.Set) (entry, error) {
	var (
		name = s.GetFirst("name")
		typ  = s.GetFirst("type")
		size = s.GetFirst("size")
		link = s.GetFirst("link")
	)
	// Normalize etnry type
	switch typ {
	case "dir", "prefix":
		typ = "dir"
	case "file", "object":
		typ = "file"
	default:
		return entry{}, fmt.Errorf("unknown entry type: %s", typ)
	}
	nsize, _ := strconv.ParseUint(size, 10, 64)
	return entry{
		Name: name,
		Type: typ,
		Size: nsize,
		Link: link,
	}, nil
}

func toOPFS(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	tmpl, err := template.ParseFS(assetFS, "assets/toopfs.html", "assets/toopfs.js", "assets/toopfs.css")
	if err != nil {
		return nil, err
	}

	// TODO:
	var d doc
	for s, err := range filterbase.NewLTSVReader(r).Iter() {
		if err != nil {
			r.Close()
			return nil, err
		}
		e, err := toEntry(s)
		if err != nil {
			return nil, err
		}
		if !e.IsDir() {
			d.Downloads = append(d.Downloads, e)
		}
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, &d); err != nil {
		return nil, err
	}
	return r.Wrap(io.NopCloser(buf)).
		PutContentType("text/html").
		Put(resource.SkipFilters, true), nil
}

func init() {
	filter.MustRegister("toopfs", toOPFS)
}
