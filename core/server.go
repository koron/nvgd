package core

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/protocol/configp"
	"github.com/koron/nvgd/resource"
)

// Server represents NVGD server.
type Server struct {
	httpd          *http.Server
	accessLog      *log.Logger
	errorLog       *log.Logger
	defaultFilters *Filters

	aliases                  aliases
	accessControlAllowOrigin string
	rootContentsFile         string

	rscSrv *resourceServer
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
	rscSrv, err := newResourceServer()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare resource/assets server: %w", err)
	}
	aliases := defaultAliases.mergeMap(c.Aliases)
	// FIXME: check conflictions between aliases and rscSrv
	// FIXME: should not be global.
	configp.Config = *c
	s := &Server{
		accessLog:                alog,
		errorLog:                 elog,
		defaultFilters:           &Filters{descs: c.DefaultFilters},
		aliases:                  aliases,
		accessControlAllowOrigin: c.AccessControlAllowOrigin,
		rootContentsFile:         c.RootContentsFile,
		rscSrv:                   rscSrv,
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
	s.accessLog.Printf("%s %s %s", req.Method, req.URL.EscapedPath(), req.URL.RawQuery)
	// "Access-Control-Allow-Origin" header enabled by configuration.
	if v := s.accessControlAllowOrigin; v != "" {
		res.Header().Set("Access-Control-Allow-Origin", v)
		if v != "*" {
			res.Header().Set("Vary", "Origin")
		}
		if req.Method == "OPTIONS" {
			res.Header().Set("Access-Control-Allow-Methods", "HEAD, GET, POST, ORIGIN")
			res.Header().Set("Access-Control-Allow-Headers", "*")
		}
	}
	// serve customized root contents
	if req.URL.Path == "/" && s.rootContentsFile != "" {
		s.serveFile(res, req, s.rootContentsFile)
		return
	}
	// serve embedded resources: favicon.ico or so.
	isServed, err := s.rscSrv.isServed(req.URL.Path)
	if err != nil {
		s.serveError(res, err)
		return
	}
	if isServed {
		s.rscSrv.serveHTTP(res, req)
		return
	}
	if err := s.serveProtocols(res, req); err != nil {
		s.serveError(res, err)
	}
}

// serveError serves an error and logs it.
func (s *Server) serveError(w http.ResponseWriter, err error) {
	msg, code := toHTTPError(err)
	http.Error(w, msg, code)
	s.errorLog.Printf("serve an error: %s", msg)
}

// serveFile serves a specified file.
func (s *Server) serveFile(w http.ResponseWriter, r *http.Request, name string) {
	f, err := os.Open(name)
	if err != nil {
		s.serveError(w, err)
		return
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		s.serveError(w, err)
		return
	}
	if fi.IsDir() {
		s.serveError(w, errors.New("root contents should not be a directory"))
		return
	}
	// Last-Modified
	modtime := fi.ModTime()
	if !modtime.IsZero() && !modtime.Equal(time.Unix(0, 0)) {
		w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	}
	// Content-Type
	ctype := mime.TypeByExtension(filepath.Ext(name))
	if ctype == "" {
		ctype = "text/plain"
	}
	w.Header().Set("Content-Type", ctype)
	// Content-Length
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	// flush header
	w.WriteHeader(http.StatusOK)
	// body
	if r.Method != "HEAD" {
		io.Copy(w, f)
	}
}

func (s *Server) isPost(p protocol.Protocol, req *http.Request) (protocol.Postable, bool) {
	if req.Method != http.MethodPost {
		return nil, false
	}
	p2, ok := p.(protocol.Postable)
	if !ok {
		return nil, false
	}
	return p2, true

}

func (s *Server) open(p protocol.Protocol, u *url.URL, req *http.Request) (*resource.Resource, error) {
	if p2, ok := s.isPost(p, req); ok {
		defer req.Body.Close()
		data := req.Body
		err := req.ParseMultipartForm(32 * 1024 * 1024)
		if err == nil {
			fh, ok := req.MultipartForm.File["file00"]
			if !ok || len(fh) < 1 {
				return nil, errors.New("no files uploaded")
			}
			f, err := fh[0].Open()
			if err != nil {
				return nil, err
			}
			defer f.Close()
			data = f
		} else if err != http.ErrNotMultipart {
			return nil, err
		}
		rsrc, err := p2.Post(u, data)
		if err != nil {
			return nil, err
		}
		return rsrc, nil
	}
	rsrc, err := p.Open(u)
	if err != nil {
		return nil, err
	}
	return rsrc, nil
}

func (s *Server) serveProtocols(res http.ResponseWriter, req *http.Request) error {
	upath := req.URL.EscapedPath()[1:]
	upath, appliedAlias := s.aliases.apply(upath)
	u, err := url.Parse(upath)
	if err != nil {
		return fmt.Errorf("failed to parse %q as URL: %w", upath, err)
	}
	u.RawQuery = req.URL.RawQuery
	p := protocol.Find(u.Scheme)
	if p == nil {
		return fmt.Errorf("not found protocol for %q", u.Scheme)
	}
	rsrc, err := s.open(p, u, req)
	if err != nil {
		return fmt.Errorf("failed to open %s; %w", upath, err)
	}
	if rsrc == nil {
		return fmt.Errorf("nil resource for %s", upath)
	}
	if redirect, ok := rsrc.String(commonconst.Redirect); ok {
		http.Redirect(res, req, "/"+redirect, http.StatusSeeOther)
		return nil
	}
	if v, ok := rsrc.Bool(commonconst.LTSV); v && ok && appliedAlias != nil {
		rewritten, err := appliedAlias.rewriteLTSV(rsrc)
		if err != nil {
			return fmt.Errorf("rewrite alias failure: %w", err)
		}
		rsrc = rewritten
	}

	// Respond to preflight requests only when the resource exists.
	if req.Method == http.MethodOptions {
		res.Header().Set("Content-Length", "0")
		res.WriteHeader(http.StatusOK)
		return nil
	}

	// Apply filters to query params.
	qp, err := qparamsParse(req.URL.RawQuery)
	if err != nil {
		rsrc.Close()
		return fmt.Errorf("failed to parse query string: %w", err)
	}
	if parsed, ok := rsrc.Strings(commonconst.ParsedKeys); ok {
		qp = qp.deleteKeys(parsed)
	}
	qp, refresh := s.splitRefresh(qp)
	qp, download := s.splitDownload(qp)
	qp, all := s.splitAll(qp)

	// Apply filters to the contents.
	r, err := s.applyFilters(qp, rsrc)
	if err != nil {
		if r != nil {
			r.Close()
		}
		return fmt.Errorf("filter error: %w", err)
	}
	if !all && !s.isSmall(rsrc) {
		r, err = s.defaultFilters.apply(s, upath, r)
		if err != nil {
			if r != nil {
				r.Close()
			}
			return fmt.Errorf("default filters for %q causes problem: %w", upath, err)
		}
	}
	defer r.Close()

	// Process special filters.

	/// Set Refresh header if required.
	if refresh > 0 {
		v := fmt.Sprintf("%d; URL=%s", refresh, req.URL.String())
		res.Header().Set("Refresh", v)
	}
	// Set Content-Disposition header if required.
	if fn, ok := r.String(resource.Filename); ok {
		res.Header().Set("Content-Disposition",
			fmt.Sprintf(`attachment; filename="%s"`, fn))
	} else if download {
		v := "attachment"
		fn := path.Base(u.Path)
		if fn != "" && fn != "." && fn != "/" {
			v = fmt.Sprintf(`attachment; filename="%s"`, fn)
		}
		res.Header().Set("Content-Disposition", v)
	}
	// Set "Content-Type" header if required.
	if ct, ok := r.String(resource.ContentType); ok {
		res.Header().Set("Content-Type", ct)
	} else if s.isHTML(qp) {
		res.Header().Set("Content-Type", "text/html")
	}

	// Output the headers and the body to ResponseWriter.
	res.WriteHeader(http.StatusOK)
	_, err = io.Copy(res, r)
	if err != nil {
		s.errorLog.Printf("failed to copy body content: %s", err)
	}
	return nil
}

func (s *Server) isSmall(r *resource.Resource) bool {
	v, ok := r.Bool(commonconst.Small)
	return ok && v
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

func (s *Server) splitDownload(q qparams) (qparams, bool) {
	downloads, others := q.split("download")
	if len(downloads) == 0 {
		return q, false
	}
	return others, true
}

func (s *Server) splitAll(q qparams) (qparams, bool) {
	all, others := q.split("all")
	if len(all) == 0 {
		return q, false
	}
	return others, true
}

func (s *Server) applyFilters(qp qparams, r *resource.Resource) (*resource.Resource, error) {
	for _, item := range qp {
		r2, err := s.applyFilter(item.name, item.value, r)
		if err != nil {
			return r, err
		}
		r = r2
	}
	return r, nil
}

func (s *Server) applyFilter(name, params string, r *resource.Resource) (*resource.Resource, error) {
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
	return item.name == "htmltable" || item.name == "indexhtml" ||
		item.name == "markdown"
}
