// simple is an example of using trdsql as a library.
package main

import (
	"log"

	"github.com/noborus/trdsql"
)

func main() {
	trd := trdsql.NewTRDSQL(
		trdsql.NewImporter(trdsql.InDelimiter(":")),
		trdsql.NewExporter(trdsql.NewWriter()),
	)
	err := trd.Exec("SELECT c1 FROM /etc/passwd")
	if err != nil {
		log.Fatal(err)
	}
}
