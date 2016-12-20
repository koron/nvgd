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

// DBParam is connection parameter for the database.
type DBParam struct {
	// Driver represents driver name for database.
	Driver string `yaml:"driver"`

	// Name represents driver-specific data source name.
	Name string `yaml:"name"`
}

// DBConfig represents configuration for DBHandler.
type DBConfig map[string]DBParam

// DBHandler is database protocol handler.
type DBHandler struct {
	Config *DBConfig

	l sync.Mutex
	h map[string]*sql.DB
}

var dbconfig DBConfig

func init() {
	protocol.MustRegister("db", &DBHandler{
		Config: &dbconfig,
		h:      make(map[string]*sql.DB),
	})
	config.RegisterProtocol("db", &dbconfig)
}

func (dbh *DBHandler) Open(u *url.URL) (io.ReadCloser, error) {
	var (
		name  = u.Host
		query = u.Path
	)
	db, err := dbh.openDB(name)
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
	return dbh.rows2ltsv(rows)
}

func (dbh *DBHandler) openDB(name string) (*sql.DB, error) {
	dbh.l.Lock()
	defer dbh.l.Unlock()
	if db, ok := dbh.h[name]; ok {
		return db, nil
	}
	p, ok := (*dbh.Config)[name]
	if !ok {
		return nil, fmt.Errorf("unknown database: %q")
	}
	return sql.Open(p.Driver, p.Name)
}

func (dbh *DBHandler) rows2ltsv(rows *sql.Rows) (io.ReadCloser, error) {
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
	for i, _ := range vals {
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
