package trdsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"unsafe"

	"github.com/noborus/trdsql"
)

// DB is a struct for forcibly rewriting the unexported fields of trdsql.DB
// using unsafe
type DB struct {
	Driver  string
	Dsn     string
	Quote   string
	MaxBulk int
	*sql.DB
	Tx          *sql.Tx
	ImportCount int
}

func toOwnDB(db *trdsql.DB) *DB {
	return (*DB)(unsafe.Pointer(db))
}

func setTrdsqlDBQuote(db *trdsql.DB, quote string) {
	toOwnDB(db).Quote = quote
}

// execTrdsql is based on a copy of trdsql.ExecContext
func execTrdsql(ctx context.Context, trd *trdsql.TRDSQL, sqlQuery string) error {
	db, err := trdsql.Connect(trd.Driver, trd.Dsn)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}

	// rewrite db.quote forcibly.
	setTrdsqlDBQuote(db, "`")

	defer func() {
		if deferErr := db.Disconnect(); deferErr != nil {
			log.Printf("disconnect: %s", deferErr)
		}
	}()

	db.Tx, err = db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin: %w", err)
	}

	if trd.Importer != nil {
		sqlQuery, err = trd.Importer.ImportContext(ctx, db, sqlQuery)
		if err != nil {
			return fmt.Errorf("import: %w", err)
		}
	}

	if trd.Exporter != nil {
		if err := trd.Exporter.ExportContext(ctx, db, sqlQuery); err != nil {
			return fmt.Errorf("export: %w", err)
		}
	}

	if err := db.Tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}
