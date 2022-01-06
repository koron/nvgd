package db

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"strings"
	"sync"

	"github.com/koron/nvgd/config"
)

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

var dbconfig Config

var (
	connPool = map[string]*conn{}
	connLock sync.Mutex
)

func init() {
	config.RegisterProtocol("db", &dbconfig)
}

func getDBParam(name string) (*Param, error) {
	p, ok := dbconfig[name]
	if !ok {
		return nil, fmt.Errorf("unknown database: %q", name)
	}
	return &p, nil
}

func (p *Param) expandName(dbname string) (string, error) {
	if dbname == "" {
		return p.Name, nil
	}
	t, err := template.New(p.Driver).Parse(p.Name)
	if err != nil {
		return "", err
	}
	t.Option("missingkey=error")
	q := map[string]string{
		"dbname": dbname,
	}
	b := &bytes.Buffer{}
	err = t.Execute(b, q)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (p *Param) openDB(dbname string) (*conn, error) {
	id := p.Name
	if p.MultipleDatabase {
		id += "--" + dbname
	} else {
		dbname = ""
	}
	connLock.Lock()
	defer connLock.Unlock()
	if c, ok := connPool[id]; ok {
		return c, nil
	}
	n, err := p.expandName(dbname)
	if err != nil {
		return nil, err
	}
	c, err := connect(p.Driver, n, p.MaxRows)
	if err != nil {
		return nil, err
	}
	connPool[id] = c
	return c, nil
}

func extractNames(u *url.URL) (name, dbname string) {
	name = u.Hostname()
	if u.User != nil {
		dbname = u.User.Username()
	}
	return name, dbname
}

func openDB(u *url.URL) (*conn, error) {
	name, dbname := extractNames(u)
	p, err := getDBParam(name)
	if err != nil {
		return nil, err
	}
	c, err := p.openDB(dbname)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func regulatePath(u *url.URL) string {
	if strings.HasPrefix(u.Path, "/") {
		return u.Path[1:]
	}
	return u.Path
}

func parseAsTables(u *url.URL) []string {
	p := regulatePath(u)
	if p == "" {
		return nil
	}
	return strings.Split(p, ",")
}
