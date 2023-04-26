// import is an example of using a customized import.
package main

import (
	"context"
	"log"

	"github.com/noborus/trdsql"
)

type importer struct{}

func (f *importer) ImportContext(ctx context.Context, db *trdsql.DB, query string) (string, error) {
	err := db.CreateTableContext(ctx, "test", []string{"a"}, []string{"text"}, true)
	if err != nil {
		return query, err
	}
	_, err = db.Tx.ExecContext(ctx, "INSERT INTO test (a) VALUES ('aaaa')")
	if err != nil {
		return query, err
	}
	return query, err
}

func (f *importer) Import(db *trdsql.DB, query string) (string, error) {
	ctx := context.Background()
	return f.ImportContext(ctx, db, query)
}

func main() {
	trdsql.EnableDebug()
	d := importer{}

	trd := trdsql.NewTRDSQL(&d, trdsql.NewExporter(trdsql.NewWriter()))
	err := trd.Exec("SELECT * FROM test")
	if err != nil {
		log.Fatal(err)
	}
}
