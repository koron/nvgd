package core

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/protocol/configp"
	"github.com/koron/nvgd/resource"
)

// Server represents NVGD server.
type Server struct {
	httpd     *http.Server
	fileSrv   http.Handler
	accessLog *log.Logger
	errorLog  *log.Logger
	filters   *Filters

	aliases aliases
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
	stripFS, err := fs.Sub(assetsFS, "assets")
	if err != nil {
		return nil, fmt.Errorf("failed to fs.Sub on core/assets: %w", err)
	}
	// FIXME: should not be global.
	configp.Config = *c
	s := &Server{
		fileSrv:   http.FileServer(http.FS(stripFS)),
		accessLog: alog,
		errorLog:  elog,
		filters:   &Filters{descs: c.Filters},
		aliases:   defaultAliases.mergeMap(c.Aliases),
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
	if req.URL.Path == "/favicon.ico" {
		s.fileSrv.ServeHTTP(res, req)
		return
	}
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

func (s *Server) serve(res http.ResponseWriter, req *http.Request) error {
	upath := req.URL.EscapedPath()[1:]
	upath, appliedAlias := s.aliases.apply(upath)
	u, err := url.Parse(upath)
	if err != nil {
		return fmt.Errorf("failed to parse %q as URL: %s", upath, err)
	}
	u.RawQuery = req.URL.RawQuery
	p := protocol.Find(u.Scheme)
	if p == nil {
		return fmt.Errorf("not found protocol for %q", u.Scheme)
	}
	rsrc, err := s.open(p, u, req)
	if err != nil {
		return fmt.Errorf("failed to open %s; %s", upath, err)
	}
	if rsrc == nil {
		return fmt.Errorf("nil resource for %s", upath)
	}
	if v, ok := rsrc.Bool(commonconst.LTSV); v && ok && appliedAlias != nil {
		rewritten, err := appliedAlias.rewriteLTSV(rsrc)
		if err != nil {
			return fmt.Errorf("rewrite alias failure: %w", err)
		}
		rsrc = rewritten
	}

	qp, err := qparamsParse(req.URL.RawQuery)
	if err != nil {
		rsrc.Close()
		return fmt.Errorf("failed to parse query string: %s", err)
	}
	if parsed, ok := rsrc.Strings(protocol.ParsedKeys); ok {
		qp = qp.deleteKeys(parsed)
	}
	qp, refresh := s.splitRefresh(qp)
	qp, download := s.splitDownload(qp)
	qp, all := s.splitAll(qp)
	r, err := s.applyFilters(qp, rsrc)
	if err != nil {
		if r != nil {
			r.Close()
		}
		return fmt.Errorf("filter error: %s", err)
	}
	if !all && !s.isSmall(rsrc) {
		r, err = s.filters.apply(s, upath, r)
		if err != nil {
			if r != nil {
				r.Close()
			}
			return fmt.Errorf("default filters for %q causes problem: %s", upath, err)
		}
	}
	defer r.Close()
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
	res.WriteHeader(http.StatusOK)
	_, err = io.Copy(res, r)
	if err != nil {
		s.errorLog.Printf("failed to copy body content: %s", err)
	}
	return nil
}

func (s *Server) isSmall(r *resource.Resource) bool {
	v, ok := r.Bool(protocol.Small)
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
