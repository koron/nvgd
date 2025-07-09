package db

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"strings"

	xlsx4db "github.com/koron/go-xlsx4db"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

type UpdateHandler struct {
}

func init() {
	protocol.MustRegister("db-update", &UpdateHandler{})
}

func (uh *UpdateHandler) Open(u *url.URL) (*resource.Resource, error) {
	name, _ := extractNames(u)
	p, err := getDBParam(name)
	if err != nil {
		return nil, err
	}
	s := u.Path
	if s == "" {
		// redirect to "/" appended URL.
		u.Path = "/"
		return resource.NewRedirect(u.String()), nil
	}
	if !strings.HasPrefix(s, "/") {
		return nil, errors.New("UpdateHandler#Open: unknown resource")
	}
	return uh.openAsset(s[1:], map[string]any{
		"name": name,
		"db":   p,
	})
}

func (uh *UpdateHandler) openAsset(s string, p map[string]any) (*resource.Resource, error) {
	if s == "" {
		s = "update.html"
	}
	f, err := assetsOpen(s)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(s, ".html") {
		// extract template.
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			return nil, err
		}
		t, err := template.New(s).Parse(string(b))
		if err != nil {
			return nil, err
		}
		bb := new(bytes.Buffer)
		err = t.Execute(bb, p)
		if err != nil {
			return nil, err
		}
		rs := resource.New(io.NopCloser(bb)).GuessContentType(s)
		return rs, nil
	}
	rs := resource.New(f).GuessContentType(s)
	return rs, nil
}

func (uh *UpdateHandler) Post(u *url.URL, r io.Reader) (*resource.Resource, error) {
	xf, err := openXLSX(r)
	if err != nil {
		return nil, fmt.Errorf("XLSX format error: %w", err)
	}
	c, err := openDB(u)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	tables := parseAsTables(u)
	err = xlsx4db.Update(c.db, xf, tables...)
	if err != nil {
		return nil, fmt.Errorf("failed to restore: %w", err)
	}
	return resource.NewString("updated successfully"), nil
}
