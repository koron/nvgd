package db

import (
	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	// PostgreSQL driver
	_ "github.com/lib/pq"
	// SQLite3 driver requires cgo, disabled as default.
	//_ "github.com/mattn/go-sqlite3"
)
