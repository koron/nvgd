package db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"github.com/koron/go-xlsx4db"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
	"github.com/tealeg/xlsx"
)

type DumpHandler struct {
}

func init() {
	protocol.MustRegister("db-dump", &DumpHandler{})
}

func (dh *DumpHandler) Open(u *url.URL) (*resource.Resource, error) {
	if p := path(u); p == "" || strings.HasPrefix(p, assetPrefix) {
		return dh.openAsset(p)
	}
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	xf := xlsx.NewFile()
	tables := parseAsTables(u)
	err = xlsx4db.Dump(xf, c.db, tables...)
	if err != nil {
		return nil, err
	}
	buf, err := saveToBuffer(xf)
	if err != nil {
		return nil, err
	}
	rs := resource.New(ioutil.NopCloser(buf))
	n, _ := extractNames(u)
	t := time.Now().Format("20060102T150405MST")
	rs.PutFilename(fmt.Sprintf("%s-%s.xlsx", n, t))
	rs.Put(protocol.Small, true)
	return rs, nil
}

func (dh *DumpHandler) openAsset(s string) (*resource.Resource, error) {
	if s == "" {
		s = "dump.html"
	}
	if strings.HasPrefix(s, assetPrefix) {
		s = s[len(assetPrefix):]
	}
	f, err := Assets.Open(s)
	if err != nil {
		return nil, err
	}
	rs := resource.New(f).GuessContentType(s)
	return rs, nil
}

func saveToBuffer(xf *xlsx.File) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := xf.Write(buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
