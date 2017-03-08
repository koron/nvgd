package db

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"text/template"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/ltsv"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

// NullReplacement replaces null value in LTSV.
var NullReplacement = "(null)"

// Param is connection parameter for the database.
type Param struct {
	// Driver represents driver name for database.
	Driver string `yaml:"driver"`

	// Name represents driver-specific data source name.
	Name string `yaml:"name"`

	// MaxRows is limitation of rows.
	MaxRows int `yaml:"max_rows"`

	// MultipleDatabase to support multiple database in an instance.
	MultipleDatabase bool `yaml:"multiple_database"`
}

// Config represents configuration for Handler.
type Config map[string]Param

// Handler is database protocol handler.
type Handler struct {
	Config *Config

	l     sync.Mutex
	conns map[string]*conn
}

var dbconfig Config

func init() {
	protocol.MustRegister("db", &Handler{
		Config: &dbconfig,
		conns:  make(map[string]*conn),
	})
	config.RegisterProtocol("db", &dbconfig)
}

// Open creates a database handler.
func (h *Handler) Open(u *url.URL) (*resource.Resource, error) {
	var (
		name, query = h.extractNameAndQuery(u)
		dbname      string
		err         error
	)
	p, ok := (*h.Config)[name]
	if !ok {
		return nil, fmt.Errorf("unknown database: %q", name)
	}
	if p.MultipleDatabase {
		dbname, query, err = h.splitDbnameAndQuery(query)
	}
	dbid := name + "--" + dbname
	c, err := h.openDB(p, dbid, dbname)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(query, "/") {
		query = query[1:]
	}
	if err := h.checkSanity(query); err != nil {
		return nil, err
	}
	rc, err := h.execQuery(c, query)
	if err != nil {
		return nil, err
	}
	return resource.New(rc), nil
}

func (h *Handler) splitDbnameAndQuery(s string) (dbname, query string, err error) {
	t := strings.SplitN(s, "/", 2)
	if len(t) != 2 || t[0] == "" {
		return "", "", fmt.Errorf("can't find database name in %q", s)
	}
	return t[0], t[1], nil
}

func (h *Handler) extractNameAndQuery(u *url.URL) (name, query string) {
	name = u.Host
	query = u.Path
	if strings.HasPrefix(query, "/") {
		query = query[1:]
	}
	return name, query
}

var reBadQuery = regexp.MustCompile(`(?i:^\s*(?:insert|update|delete|create|drop|alter|truncate|prepare|execute))`)

func (h *Handler) checkSanity(q string) error {
	// FIXME: too simple, should do more.
	if reBadQuery.MatchString(q) {
		return errors.New("including invalid keywords")
	}
	return nil
}

// execQuery executes a query in a transaction which will be rollbacked.
func (h *Handler) execQuery(c *conn, q string) (io.ReadCloser, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	rows, err := tx.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return h.rows2ltsv(rows, c.maxRows)
}

func (h *Handler) openDB(p Param, dbid, dbname string) (*conn, error) {
	h.l.Lock()
	defer h.l.Unlock()
	if c, ok := h.conns[dbid]; ok {
		return c, nil
	}
	name, err := h.expandName(p.Driver, p.Name, dbname)
	if err != nil {
		return nil, err
	}
	c, err := connect(p.Driver, name, p.MaxRows)
	if err != nil {
		return nil, err
	}
	h.conns[dbid] = c
	return c, nil
}

func (h *Handler) expandName(driver, name, dbname string) (string, error) {
	if dbname == "" {
		return name, nil
	}
	t, err := template.New(driver).Parse(name)
	t.Option("missingkey=error")
	if err != nil {
		return "", err
	}
	p := map[string]string {
		"dbname": dbname,
	}
	b := &bytes.Buffer{}
	err = t.Execute(b, p)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (h *Handler) rows2ltsv(rows *sql.Rows, maxRows int) (io.ReadCloser, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var (
		buf = &bytes.Buffer{}
		w   = ltsv.NewWriter(buf, cols...)
		n   = len(cols)
	)

	vals := make([]interface{}, n)
	for i := range vals {
		vals[i] = new(sql.NullString)
	}
	strs := make([]string, n)

	nrow := 0
	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return nil, err
		}
		for i, v := range vals {
			ns := v.(*sql.NullString)
			if ns.Valid {
				strs[i] = ns.String
			} else {
				strs[i] = NullReplacement
			}
		}
		w.Write(strs...)
		nrow++
		if nrow >= maxRows {
			break
		}
	}
	return ioutil.NopCloser(buf), nil
}
