package protocol_test

import (
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/koron/nvgd/internal/protocoltest"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

type dummyPoster struct{}

func (*dummyPoster) Open(u *url.URL) (*resource.Resource, error) {
	return resource.NewString("_postable_dummy_open_"), nil
}

func (*dummyPoster) Post(u *url.URL, r io.Reader) (*resource.Resource, error) {
	var body string
	b, err := io.ReadAll(r)
	if err != nil {
		body = fmt.Sprintf("read error: %s", err)
	} else {
		body = string(b)
	}
	msg := fmt.Sprintf("_postable_dummy_post_: body=%s", body)
	return resource.NewString(msg), nil
}

func init() {
	protocol.MustRegister("dummy", protocol.ProtocolFunc(dummyProtocol))
	protocol.MustRegister("dummypost", &dummyPoster{})
}

func dummyProtocol(u *url.URL) (*resource.Resource, error) {
	return resource.NewString("_dummy_content_"), nil
}

func TestRegister(t *testing.T) {
	protocol.MustRegister("__register__", protocol.ProtocolFunc(dummyProtocol))
	protocoltest.CheckRegistered(t, "__register__", protocol.ProtocolFunc(dummyProtocol))
}

func TestDuplicate(t *testing.T) {
	protocoltest.CheckRegistered(t, "dummy", protocol.ProtocolFunc(dummyProtocol))
	defer func() {
		err := recover().(error)
		if err == nil {
			t.Fatal("no duplications, unexpectedly")
		}
		got := err.Error()
		want := `duplicated protocol name "dummy"`
		if got != want {
			t.Fatalf("unexpected error message:\nwant=%s\n got=%s", want, got)
		}
	}()
	protocol.MustRegister("dummy", protocol.ProtocolFunc(dummyProtocol))
}

func TestOpen(t *testing.T) {
	got := protocoltest.OpenString(t, "dummy://")
	want := "_dummy_content_"
	if got != want {
		t.Errorf("unmatch content:\nwant=%s\n got=%s", want, got)
	}
}

func TestOpenFail(t *testing.T) {
	protocoltest.OpenFail(t, "unexist://")
}

func TestPost(t *testing.T) {
	rsrc := protocoltest.Post(t, "dummypost://", "_post_content_")
	got := protocoltest.ReadAllString(t, rsrc)
	want := `_postable_dummy_post_: body=_post_content_`
	if got != want {
		t.Errorf("unmatch response:\nwant=%s\n got=%s", want, got)
	}
}

func TestPostMultipart(t *testing.T) {
	rsrc := protocoltest.Post(t, "dummypost://", map[string]string{
		"file00": "_file00_content_",
		"file01": "_file01_content_",
		"field0": "_field0_content_",
		"field1": "_field1_content_",
		"field2": "_field2_content_",
	})
	got := protocoltest.ReadAllString(t, rsrc)
	want := `_postable_dummy_post_: body=_file00_content_`
	if got != want {
		t.Errorf("unmatch response:\nwant=%s\n got=%s", want, got)
	}
}
