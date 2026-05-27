package db

import (
	"database/sql"
	"database/sql/driver"
	"sync"
	"sync/atomic"
	"testing"
)

type closeTestDriver struct{}

type closeTestConn struct{}

func (d *closeTestDriver) Open(name string) (driver.Conn, error) {
	return &closeTestConn{}, nil
}

func (c *closeTestConn) Prepare(query string) (driver.Stmt, error) {
	return nil, nil
}

func (c *closeTestConn) Close() error {
	return nil
}

func (c *closeTestConn) Begin() (driver.Tx, error) {
	return nil, nil
}

func init() {
	sql.Register("close_test", &closeTestDriver{})
}

func TestConnCloseConcurrent(t *testing.T) {
	db, err := sql.Open("close_test", "")
	if err != nil {
		t.Fatal(err)
	}
	c := &conn{id: "test", db: db}

	var panicked atomic.Bool
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					panicked.Store(true)
				}
			}()
			c.Close()
		}()
	}
	wg.Wait()
	if panicked.Load() {
		t.Fatal("panic during concurrent Close")
	}
}
