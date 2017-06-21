package db

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
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
		// TODO: redirect
		return nil, errors.New("UpdateHandler#Open: try to append '/' to URL")
	}
	if !strings.HasPrefix(s, "/") {
		return nil, errors.New("UpdateHandler#Open: unknown resource")
	}
	return uh.openAsset(s[1:], map[string]interface{}{
		"name": name,
		"db":   p,
	})
	return nil, errors.New("UpdateHandler#Open: not implemented yet")
}

func (uh *UpdateHandler) openAsset(s string, p map[string]interface{}) (*resource.Resource, error) {
	if s == "" {
		s = "update.html"
	}
	f, err := Assets.Open(s)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(s, ".html") {
		// extract template.
		defer f.Close()
		b, err := ioutil.ReadAll(f)
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
		rs := resource.New(ioutil.NopCloser(bb)).GuessContentType(s)
		return rs, nil
	}
	rs := resource.New(f).GuessContentType(s)
	return rs, nil
}

func (uh *UpdateHandler) Post(u *url.URL, r io.Reader) (*resource.Resource, error) {
	xf, err := openXLSX(r)
	if err != nil {
		return nil, err
	}
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	tables := parseAsTables(u)
	err = xlsx4db.Update(c.db, xf, tables...)
	if err != nil {
		return nil, err
	}
	return resource.NewString("updated successfully"), nil
}
