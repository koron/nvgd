// Package toopfs provides UI to download files to the OPFS.
package toopfs

import (
	"bytes"
	"embed"
	"html/template"
	"io"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/fileentry"
	"github.com/koron/nvgd/internal/filterbase"
	"github.com/koron/nvgd/resource"
)

//go:embed assets
var embedFS embed.FS

var assetFS = devfs.New(embedFS, "filter/toopfs", "")

type doc struct {
	Downloads []fileentry.Entry
}

func toOPFS(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	tmpl, err := template.ParseFS(assetFS, "assets/toopfs.html", "assets/toopfs.js", "assets/toopfs.css")
	if err != nil {
		return nil, err
	}

	var d doc
	for s, err := range filterbase.NewLTSVReader(r).Iter() {
		if err != nil {
			return nil, err
		}
		entry, err := fileentry.ParseLTSV(s)
		if err != nil {
			return nil, err
		}
		d.Downloads = append(d.Downloads, entry)
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
