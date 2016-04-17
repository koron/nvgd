package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/protocol"
)

// Server represents NVGD server.
type Server struct {
	httpd *http.Server
	log   *log.Logger
}

// New creates a server instance.
func New(c *Config) (*Server, error) {
	logger, err := c.logger()
	if err != nil {
		return nil, err
	}
	s := &Server{
		log: logger,
	}
	s.httpd = &http.Server{
		Addr:    c.addr(),
		Handler: s,
	}
	return s, nil
}

// Run runs NGVD server.
func (s *Server) Run() error {
	return s.httpd.ListenAndServe()
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	u, err := url.Parse(path)
	if err != nil {
		// TODO: error response
		s.log.Printf("failed to parse %q as URL: %s", path, err)
		return
	}
	p := protocol.Find(u.Scheme)
	if p == nil {
		// TODO: error response
		s.log.Printf("not found protocol for %q", u.Scheme)
		return
	}
	r, err := p.Open(u)
	if err != nil {
		// TODO: error response
		s.log.Printf("failed to open: %s", path)
		return
	}
	r, err = s.applyFilters(u.RawQuery, r)
	if err != nil {
		if r != nil {
			r.Close()
		}
		// TODO: error response
		s.log.Printf("filter error: %s", err)
		return
	}
	defer r.Close()
	// TODO: better log
	s.log.Printf("%s %s", req.Method, req.URL.Path)
	res.WriteHeader(http.StatusOK)
	_, err = io.Copy(res, r)
}

func (s *Server) applyFilters(q string, r io.ReadCloser) (io.ReadCloser, error) {
	for q != "" {
		k := q
		if i := strings.Index(k, "&"); i >= 0 {
			k, q = k[:i], k[i+1:]
		} else {
			q = ""
		}
		if k == "" {
			continue
		}
		v := ""
		if i := strings.Index(k, "="); i >= 0 {
			k, v = k[:i], k[i+1:]
		}
		k, err := url.QueryUnescape(k)
		if err != nil {
			return r, err
		}
		v, err = url.QueryUnescape(v)
		if err != nil {
			return r, err
		}
		r2, err := s.applyFilter(k, v, r)
		if err != nil {
			return r, err
		}
		r = r2
	}
	return r, nil
}

func (s *Server) applyFilter(name, params string, r io.ReadCloser) (io.ReadCloser, error) {
	f := filter.Find(name)
	if f == nil {
		return nil, fmt.Errorf("not found filter: %s", name)
	}
	p, err := s.parseParams(params)
	if err != nil {
		return nil, err
	}
	return f.Filter(r, p)
}

func (s *Server) parseParams(q string) (map[string]string, error) {
	p := map[string]string{}
	for q != "" {
		k := q
		if i := strings.Index(k, ";"); i >= 0 {
			k, q = k[:i], k[i+1:]
		} else {
			q = ""
		}
		if k == "" {
			continue
		}
		v := ""
		if i := strings.Index(k, ":"); i >= 0 {
			k, v = k[:i], k[i+1:]
		}
		p[k] = v
	}
	return p, nil
}
