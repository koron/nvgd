package core

import (
	"log"
	"net/http"
)

// Server represents NVGD server.
type Server struct {
	httpd *http.Server
}

// New creates a server instance.
func New(c *Config) (*Server, error) {
	s := &Server{}
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
	log.Printf("%s %s", req.Method, req.URL.Path)
	// TODO:
}
