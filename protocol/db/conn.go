package db

import (
	"database/sql"
	"sync"
)

const defaultMaxRows = 100

var (
	connPool = map[string]*conn{}
	connLock sync.Mutex
)

type conn struct {
	id      string
	driver  string
	db      *sql.DB
	maxRows int
}

func connect(driver, name string, maxRows int) (*conn, error) {
	id := driver + "--" + name
	connLock.Lock()
	defer connLock.Unlock()
	if c, ok := connPool[id]; ok {
		return c, nil
	}
	db, err := sql.Open(driver, name)
	if err != nil {
		return nil, err
	}
	if maxRows <= 0 {
		maxRows = defaultMaxRows
	}
	c := &conn{
		id:      id,
		driver:  driver,
		db:      db,
		maxRows: maxRows,
	}
	connPool[id] = c
	return c, nil
}

// Close closes underlying connection.
func (c *conn) Close() error {
	if c.db == nil {
		return nil
	}
	connLock.Lock()
	defer connLock.Unlock()
	delete(connPool, c.id)
	err := c.db.Close()
	c.db = nil
	return err
}
