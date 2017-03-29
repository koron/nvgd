package aws

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func newProxyFunc(proxy string) func(*http.Request) (*url.URL, error) {
	return func(*http.Request) (*url.URL, error) {
		// copy from ProxyFromEnvironment() in net/http/transport.go
		proxyURL, err := url.Parse(proxy)
		if err != nil || !strings.HasPrefix(proxyURL.Scheme, "http") {
			// proxy was bogus. Try prepending "http://" to it and
			// see if that parses correctly. If not, we fall
			// through and complain about the original one.
			if proxyURL, err := url.Parse("http://" + proxy); err == nil {
				return proxyURL, nil
			}
		}
		if err != nil {
			return nil, fmt.Errorf("invalid proxy address %q: %v", proxy, err)
		}
		return proxyURL, nil
	}
}

func NewProxyClient(proxy string) *http.Client {
	t := *(http.DefaultTransport.(*http.Transport))
	t.Proxy = newProxyFunc(proxy)
	c := *http.DefaultClient
	c.Transport = &t
	return &c
}
