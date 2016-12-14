package core

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/protocol"
)

// Server represents NVGD server.
type Server struct {
	httpd     *http.Server
	accessLog *log.Logger
	errorLog  *log.Logger
}

// New creates a server instance.
func New(c *config.Config) (*Server, error) {
	alog, err := c.AccessLog()
	if err != nil {
		return nil, err
	}
	elog, err := c.ErrorLog()
	if err != nil {
		return nil, err
	}
	s := &Server{
		accessLog: alog,
		errorLog:  elog,
	}
	s.httpd = &http.Server{
		Addr:    c.Addr,
		Handler: s,
	}
	s.errorLog.Printf("start to listening on %s", c.Addr)
	return s, nil
}

// Run runs NGVD server.
func (s *Server) Run() error {
	return s.httpd.ListenAndServe()
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	s.accessLog.Printf("%s %s %s", req.Method, req.URL.Path, req.URL.RawQuery)
	if err := s.serve(res, req); err != nil {
		// TODO: log an error.
		if herr, ok := err.(httpError); ok {
			res.WriteHeader(herr.statusCode())
			res.Write(([]byte)(herr.body()))
			return
		}
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(([]byte)(err.Error()))
		return
	}
}

func (s *Server) serve(res http.ResponseWriter, req *http.Request) error {
	path := req.URL.Path[1:]
	// rewrite "/files/" to "/file:///" for compatibility with koron/night.
	const filesPrefix = "files/"
	if strings.HasPrefix(path, filesPrefix) {
		path = "file:///" + path[len(filesPrefix):]
	}
	u, err := url.Parse(path)
	if err != nil {
		return fmt.Errorf("failed to parse %q as URL: %s", path, err)
	}
	p := protocol.Find(u.Scheme)
	if p == nil {
		return fmt.Errorf("not found protocol for %q", u.Scheme)
	}
	r, err := p.Open(u)
	if err != nil {
		return fmt.Errorf("failed to open %s; %s", path, err)
	}
	qp, err := qparamsParse(req.URL.RawQuery)
	if err != nil {
		return fmt.Errorf("failed to parse query string: %s", err)
	}
	qp, refresh := s.splitRefresh(qp)
	r, err = s.applyFilters(qp, r)
	if err != nil {
		if r != nil {
			r.Close()
		}
		return fmt.Errorf("filter error: %s", err)
	}
	defer r.Close()
	if refresh > 0 {
		v := fmt.Sprintf("%d; URL=%s", refresh, req.URL.String())
		res.Header().Set("Refresh", v)
	}
	if s.isHTML(qp) {
		res.Header().Set("Content-Type", "text/html")
	}
	res.WriteHeader(http.StatusOK)
	_, err = io.Copy(res, r)
	if err != nil {
		s.errorLog.Printf("failed to copy body content: %s", err)
	}
	return nil
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

func (s *Server) isHTML(qp qparams) bool {
	if len(qp) == 0 {
		return false
	}
	item := qp[len(qp)-1]
	return item.name == "htmltable"
}
