package trdsql

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

var csv1 = `id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`

func TestCSVAccumulation(t *testing.T) {
	t.Run("COUNT", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(*) FROM t",
			"ih": "true",
		}, csv1, "4\n")
	})
	t.Run("SUM", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT SUM(price) FROM t",
			"ih": "true",
		}, csv1, "2149\n")
	})
}

func TestCSVBasicOptions(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q": "SELECT * FROM t",
		}, csv1, `id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`)
	})

	t.Run("output header", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT * FROM t",
			"oh": "true",
		}, csv1, `c1,c2,c3
id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`)
	})

	t.Run("input header", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT * FROM t",
			"oh": "true",
			"ih": "true",
		}, csv1, `id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`)
	})
}

func TestLTSVOutput(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":    "SELECT * FROM t",
			"ih":   "true",
			"ofmt": "LTSV",
		}, csv1, `id:1	name:foo	price:500
id:2	name:bar	price:400
id:3	name:baz	price:250
id:4	name:qux	price:999
`)
	})
	t.Run("lower case", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":    "SELECT * FROM t",
			"ih":   "true",
			"ofmt": "ltsv",
		}, csv1, `id:1	name:foo	price:500
id:2	name:bar	price:400
id:3	name:baz	price:250
id:4	name:qux	price:999
`)
	})
}

var tsv1 = `id	name	price
1	foo	500
2	bar	400
3	baz	250
4	qux	999
`

func TestTSVInput(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":    "SELECT * FROM t",
			"ifmt": "TSV",
		}, tsv1, `id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`)
	})

	t.Run("output header", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":    "SELECT * FROM t",
			"ifmt": "tsv",
			"oh":   "true",
		}, tsv1, `c1,c2,c3
id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`)
	})

	t.Run("input header", func(t *testing.T) {
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":    "SELECT * FROM t",
			"ifmt": "tsv",
			"oh":   "true",
			"ih":   "true",
		}, tsv1, `id,name,price
1,foo,500
2,bar,400
3,baz,250
4,qux,999
`)
	})
}

func TestAttachDatabaseFail(t *testing.T) {
	name := filepath.Join(t.TempDir(), "test.db")
	filtertest.Fail(t, trdsqlFilter, filter.Params{
		"q": fmt.Sprintf("ATTACH DATABASE '%s' AS extdb", name),
	}, "N/A", "trdsql error: export: too many attached databases - max 0: ")
}

func TestManyColumns(t *testing.T) {
	genRecords := func(prefix string, n int) []string {
		recs := make([]string, n)
		for i := range n {
			recs[i] = prefix + strconv.Itoa(i)
		}
		return recs
	}

	genCsv := func(cols, rows int) string {
		t.Helper()
		bb := &bytes.Buffer{}
		w := csv.NewWriter(bb)
		err := w.Write(genRecords("c_", cols))
		if err != nil {
			t.Fatalf("failed to generate CSV header: %s", err)
		}
		for r := range rows {
			err := w.Write(genRecords("r"+strconv.Itoa(r)+"_", cols))
			if err != nil {
				t.Fatalf("failed to generate CSV body: %s", err)
			}
		}
		w.Flush()
		return bb.String()
	}

	t.Run("500", func(t *testing.T) {
		csv := genCsv(500, 10)
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_0) FROM t",
			"ih": "true",
		}, csv, "10\n")
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_499) FROM t",
			"ih": "true",
		}, csv, "10\n")
	})

	t.Run("1000", func(t *testing.T) {
		csv := genCsv(1000, 10)
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_0) FROM t",
			"ih": "true",
		}, csv, "10\n")
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_999) FROM t",
			"ih": "true",
		}, csv, "10\n")
	})

	t.Run("1500", func(t *testing.T) {
		csv := genCsv(1500, 10)
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_0) FROM t",
			"ih": "true",
		}, csv, "10\n")
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_1499) FROM t",
			"ih": "true",
		}, csv, "10\n")
	})

	t.Run("2000", func(t *testing.T) {
		csv := genCsv(2000, 10)
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_0) FROM t",
			"ih": "true",
		}, csv, "10\n")
		filtertest.Check(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_1999) FROM t",
			"ih": "true",
		}, csv, "10\n")
	})

	t.Run("2001_fail", func(t *testing.T) {
		csv := genCsv(2001, 10)
		filtertest.Fail(t, trdsqlFilter, filter.Params{
			"q":  "SELECT COUNT(c_2000) FROM t",
			"ih": "true",
		}, csv, "trdsql error: import: too many columns on t: ")
	})
}
