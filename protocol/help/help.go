// Package help provides help protocol for NVGD.
package help

import (
	"errors"
	"io/fs"
	"net/url"
	"strings"

	"github.com/koron/nvgd/doc"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func init() {
	protocol.MustRegister("help", protocol.ProtocolFunc(Serve))
}

// Text is default content of help.
var Text string

func Serve(u *url.URL) (*resource.Resource, error) {
	if u.Path == "" || strings.HasSuffix(u.Path, "/doc/") {
		u.Path = "/"
		return resource.NewRedirect(u.String()), nil
	}
	if u.Path == "/" {
		return resource.NewString(Text).PutFilename("README.md"), nil
	}
	reqPath := strings.TrimPrefix(u.Path, "/doc/")
	f, err := doc.FS.Open(reqPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			u.Path = "/"
			return resource.NewRedirect(u.String()), nil
		}
		return nil, err
	}
	return resource.New(f).PutFilename(reqPath), nil
}
