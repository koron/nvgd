package aws

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
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

// NewProxyClient creates a proxied client for AWS.
func NewProxyClient(proxy string) *http.Client {
	c := *http.DefaultClient
	c.Transport = &http.Transport{
		Proxy: newProxyFunc(proxy),
		// copy from "net/http".DefaultTransport
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return &c
}
