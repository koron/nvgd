package core

import (
	"log"
	"net/http"
	"net/url"

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

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	u, err := url.Parse(path)
	if err != nil {
		s.log.Printf("failed to parse %q as URL: %s", path, err)
		// TODO: error response
		return
	}
	p := protocol.Find(u.Scheme)
	if p == nil {
		s.log.Printf("not found protocol for %q", u.Scheme)
		// TODO: error response
		return
	}
	r, err := p.Open(u.Path)
	// TODO: fix me.
	defer r.Close()
	s.log.Printf("%s %s", req.Method, req.URL.Path)
}
