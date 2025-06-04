/*
Package protocoltest provides help utilities to implement protocol tests.
*/
package protocoltest

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

func CheckRegistered(t *testing.T, name string, want protocol.Protocol) {
	t.Helper()
	got := protocol.Find(name)
	if got == nil {
		t.Fatalf("%s: the protocol is not registered", name)
	}
	gotp := reflect.ValueOf(got).Pointer()
	wantp := reflect.ValueOf(want).Pointer()
	if gotp != wantp {
		t.Errorf("%s: unexpected protocol: wantP=0x%x got=0x%x", name, wantp, gotp)
	}
}

func GetRegistered[T any](t *testing.T, name string) T {
	t.Helper()
	var zero T
	got := protocol.Find(name)
	if got == nil {
		t.Fatalf("%s: the protocol is not registered", name)
	}
	p, ok := got.(T)
	if !ok {
		t.Fatalf("%s: registered protocol is not %T (got %T)", name, zero, got)
	}
	return p
}

func open(t *testing.T, protocolURL string, req *http.Request) *resource.Resource {
	t.Helper()
	u, err := url.Parse(protocolURL)
	if err != nil {
		t.Fatalf("failed to parse URL %s: %s", protocolURL, err)
	}
	r, err := protocol.Open(u, req)
	if err != nil {
		t.Fatalf("protocol.Open failed %q: %s", u.String(), err)
	}
	return r
}

func Open(t *testing.T, protocolURL string) *resource.Resource {
	return open(t, protocolURL, nil)
}

func OpenFail(t *testing.T, protocolURL string) error {
	t.Helper()
	u, err := url.Parse(protocolURL)
	if err != nil {
		t.Fatalf("failed to parse URL %s: %s", protocolURL, err)
	}
	_, err = protocol.Open(u, nil)
	if err == nil {
		t.Fatalf("unexpected success. expected failure: %s", u.String())
	}
	return err
}

func ReadAllString(t *testing.T, rsrc *resource.Resource) string {
	t.Helper()
	if rsrc == nil {
		t.Fatal("no resource found")
	}
	b, err := io.ReadAll(rsrc)
	if err != nil {
		t.Fatalf("failed to read the resource: %s", err)
	}
	return string(b)
}

func OpenString(t *testing.T, protocolURL string) string {
	t.Helper()
	rsrc := Open(t, protocolURL)
	defer rsrc.Close()
	return ReadAllString(t, rsrc)
}

func CheckRedirect(t *testing.T, rsrc *resource.Resource, redirectPath string) {
	t.Helper()
	got, ok := rsrc.Options[commonconst.Redirect]
	if !ok {
		t.Fatal("no redirect in the resource")
	}
	if got != redirectPath {
		t.Fatalf("unmatch redirect path: want=%q got=%q", redirectPath, got)
	}
	gotmsg := ReadAllString(t, rsrc)
	wantmsg := "redirect to: " + redirectPath
	if gotmsg != wantmsg {
		t.Fatalf("unmatch redirect message:\nwant=%q\ngot=%q", wantmsg, gotmsg)
	}
}

func CheckNotRedirect(t *testing.T, rsrc *resource.Resource) {
	t.Helper()
	if got, ok := rsrc.Options[commonconst.Redirect]; ok {
		t.Errorf("unexpected redirect to: %s", got)
	}
}

func Post(t *testing.T, protocolURL string, contents any) *resource.Resource {
	t.Helper()
	var body io.Reader
	var contentType string
	switch v := contents.(type) {
	case string:
		body = strings.NewReader(v)
	case map[string]string:
		body, contentType = multipartBytes(t, v)
	case []byte:
		body = bytes.NewReader(v)
	case io.Reader:
		body = v
	default:
		t.Fatalf("unsupported contents type: %T", contents)
	}
	req, err := http.NewRequest("POST", protocolURL, body)
	if err != nil {
		t.Fatalf("failed to create a post request: %s", err)
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return open(t, protocolURL, req)
}

func multipartBytes(t *testing.T, values map[string]string) (io.Reader, string) {
	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	defer w.Close()
	for k, v := range values {
		if strings.HasPrefix(k, "file") {
			p, err := w.CreateFormFile(k, k)
			if err != nil {
				t.Fatalf("failed to create a file: name=%s value=%s: %s", k, v, err)
			}
			io.WriteString(p, v) // bytes.Buffer never fail.
			continue
		}
		if err := w.WriteField(k, v); err != nil {
			t.Fatalf("failed to write field: name=%s value=%s: %s", k, v, err)
		}
	}
	return b, w.FormDataContentType()
}
