package trdsql

import (
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
