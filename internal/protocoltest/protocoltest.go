package protocoltest

import (
	"io"
	"net/url"
	"reflect"
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
		t.Fatalf("%s: registered protocol is not %T", name, zero)
	}
	return p
}

func Open(t *testing.T, protocolUrl string) *resource.Resource {
	t.Helper()
	u, err := url.Parse(protocolUrl)
	if err != nil {
		t.Fatalf("failed to parse URL %s: %s", protocolUrl, err)
	}
	p := protocol.Find(u.Scheme)
	if p == nil {
		t.Fatalf("protocol %q is not found", u.Scheme)
	}
	r, err := p.Open(u)
	if err != nil {
		t.Fatalf("protocol.Open failed %s: %s", u.String(), err)
	}
	return r
}

func OpenFail(t *testing.T, protocolUrl string) error {
	t.Helper()
	u, err := url.Parse(protocolUrl)
	if err != nil {
		t.Fatalf("failed to parse URL %s: %s", protocolUrl, err)
	}
	p := protocol.Find(u.Scheme)
	if p == nil {
		t.Fatalf("protocol %q is not found", u.Scheme)
	}
	_, err = p.Open(u)
	if err == nil {
		t.Fatal("unexpected success. expected failure: %s", u.String())
	}
	return err
}

func ReadAllString(t *testing.T, rsrc *resource.Resource) string {
	t.Helper()
	if rsrc == nil {
		t.Fatalf("no resource found", rsrc)
	}
	b, err := io.ReadAll(rsrc)
	if err != nil {
		t.Fatalf("failed to read the resource: %s", err)
	}
	return string(b)
}

func OpenString(t *testing.T, protocolUrl string) string {
	t.Helper()
	rsrc := Open(t, protocolUrl)
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
