// Package embedresource provides adapter between fs.FS and resource.Resource
package embedresource

import (
	"errors"
	"io/fs"
	"net/url"
	"path"

	"github.com/koron/nvgd/resource"
)

type EmbedResource struct {
	fs       fs.FS
	prefix   string
	fallback string
}

func New(fs fs.FS, opts ...Option) *EmbedResource {
	res := &EmbedResource{
		fs:       fs,
		prefix:   "assets",
		fallback: "index.html",
	}
	for _, opt := range opts {
		opt.apply(res)
	}
	return res
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
		// fall back
		if res.fallback != "" && errors.Is(err, fs.ErrNotExist) {
			reqPath = res.fallback
			f, err = res.fs.Open(path.Join(res.prefix, reqPath))
			if err == nil {
				return resource.New(f).GuessContentType(reqPath), nil
			}
		}
		return nil, err
	}
	return resource.New(f).GuessContentType(reqPath), nil
}
