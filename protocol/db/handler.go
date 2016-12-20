package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"sync"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/ltsv"
	"github.com/koron/nvgd/protocol"
)

// Param is connection parameter for the database.
type Param struct {
	// Driver represents driver name for database.
	Driver string `yaml:"driver"`

	// Name represents driver-specific data source name.
	Name string `yaml:"name"`
}

// Config represents configuration for Handler.
type Config map[string]Param

// Handler is database protocol handler.
type Handler struct {
	Config *Config

	l         sync.Mutex
	databases map[string]*sql.DB
}

var dbconfig Config

func init() {
	protocol.MustRegister("db", &Handler{
		Config:    &dbconfig,
		databases: make(map[string]*sql.DB),
	})
	config.RegisterProtocol("db", &dbconfig)
}

// Open creates a database handler.
func (h *Handler) Open(u *url.URL) (io.ReadCloser, error) {
	var (
		name  = u.Host
		query = u.Path
	)
	db, err := h.openDB(name)
	if err != nil {
		return nil, err
	}
	// TODO: sanitize query!
	if strings.HasPrefix(query, "/") {
		query = query[1:]
	}
	fmt.Printf("query=%s\n", query)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return h.rows2ltsv(rows)
}

func (h *Handler) openDB(name string) (*sql.DB, error) {
	h.l.Lock()
	defer h.l.Unlock()
	if db, ok := h.databases[name]; ok {
		return db, nil
	}
	p, ok := (*h.Config)[name]
	if !ok {
		return nil, fmt.Errorf("unknown database: %q", name)
	}
	db, err := sql.Open(p.Driver, p.Name)
	if err != nil {
		return nil, err
	}
	h.databases[name] = db
	return db, nil
}

func (h *Handler) rows2ltsv(rows *sql.Rows) (io.ReadCloser, error) {
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
		vals[i] = new(string)
	}
	strs := make([]string, n)

	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return nil, err
		}
		for i, v := range vals {
			strs[i] = *v.(*string)
		}
		w.Write(strs...)
	}
	return ioutil.NopCloser(buf), nil
}
