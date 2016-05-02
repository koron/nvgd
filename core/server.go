package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
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
	qp, err := qparamsParse(req.URL.RawQuery)
	if err != nil {
		s.log.Printf("failed to parse query string: %s", err)
		return
	}
	qp, refresh := s.splitRefresh(qp)
	r, err = s.applyFilters(qp, r)
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
	s.log.Printf("%s %s %s", req.Method, req.URL.Path, req.URL.RawQuery)
	if refresh > 0 {
		v := fmt.Sprintf("%d; URL=%s", refresh, req.URL.String())
		res.Header().Set("Refresh", v)
	}
	res.WriteHeader(http.StatusOK)
	_, err = io.Copy(res, r)
}

func (s *Server) splitRefresh(q qparams) (qparams, int) {
	refreshes, others := q.split("refresh")
	if len(refreshes) == 0 {
		return q, 0
	}
	n, err := strconv.Atoi(refreshes[0].value)
	if err != nil && n < 0 {
		n = 0
	}
	return others, n
}

func (s *Server) applyFilters(qp qparams, r io.ReadCloser) (io.ReadCloser, error) {
	for _, item := range qp {
		r2, err := s.applyFilter(item.name, item.value, r)
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
	return f(r, p)
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
