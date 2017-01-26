package db

import "database/sql"

const defaultMaxRows = 100

type conn struct {
	db      *sql.DB
	maxRows int
}

func connect(driver, name string, maxRows int) (*conn, error) {
	db, err := sql.Open(driver, name)
	if err != nil {
		return nil, err
	}
	if maxRows <= 0 {
		maxRows = defaultMaxRows
	}
	return &conn{
		db:      db,
		maxRows: maxRows,
	}, nil
}
