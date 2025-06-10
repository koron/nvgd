// Package opfs provides a UI operates with OPFS.
package opfs

import (
	"embed"
	"io/fs"
	"net/url"
	"sync"

	"github.com/koron/nvgd/internal/devfs"
	"github.com/koron/nvgd/internal/templateresource"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

//go:embed assets
var assetFS embed.FS

var getRsrc = sync.OnceValues(func() (*templateresource.TemplateResource, error) {
	fsys, err := fs.Sub(devfs.New(assetFS, "protocol/opfs", ""), "assets")
	if err != nil {
		return nil, err
	}
	return templateresource.New(fsys)
})

func init() {
	protocol.MustRegister("opfs", protocol.ProtocolFunc(open))
}

func open(u *url.URL) (*resource.Resource, error) {
	tmplRsrc, err := getRsrc()
	if err != nil {
		return nil, err
	}
	r, err := tmplRsrc.Open(u)
	if err != nil {
		return nil, err
	}
	r.Put(resource.SkipFilters, true)
	return r, nil
}
