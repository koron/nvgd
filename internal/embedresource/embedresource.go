// Package embedresource provides adapter between embed.FS and resource.Resource
package embedresource

import (
	"embed"
	"net/url"
	"path"

	"github.com/koron/nvgd/resource"
)

type EmbedResource struct {
	fs embed.FS
	prefix string
}

func New(fs embed.FS) *EmbedResource {
	return &EmbedResource{
		fs:     fs,
		prefix: "assets",
	}
}

func (res *EmbedResource) Open(u *url.URL) (*resource.Resource, error) {
	if u.Path == "" {
		u.Path = "/"
		return resource.NewRedirect(u.String()), nil
	}
	reqPath := u.Path
	if reqPath == "/" {
		reqPath = "index.html"
	}
	f, err := res.fs.Open(path.Join(res.prefix, reqPath))
	if err != nil {
		return nil, err
	}
	return resource.New(f).GuessContentType(reqPath), nil
}
