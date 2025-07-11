package texttable

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestBasic(t *testing.T) {
	const in = `
id:1	name:foo	price:123
id:2	name:bar	price:456
id:3	name:baz	price:789
`
	const want = `
+----+------+-------+
| ID | NAME | PRICE |
+----+------+-------+
|  1 | foo  |   123 |
|  2 | bar  |   456 |
|  3 | baz  |   789 |
+----+------+-------+
`
	filtertest.Check(t, filterFunc, filter.Params{}, in[1:], want[1:])
}

func TestOthers(t *testing.T) {
	const in = `
id:1	name:foo	price:123
id:2	name:bar	price:456	bar:abc
id:3	name:baz	price:789	baz:xyz	qu:x
`
	const want = `
+----+------+-------+---------------+
| ID | NAME | PRICE |   (OTHERS)    |
+----+------+-------+---------------+
|  1 | foo  |   123 | (none)        |
|  2 | bar  |   456 | bar:abc       |
|  3 | baz  |   789 | baz:xyz, qu:x |
+----+------+-------+---------------+
`
	filtertest.Check(t, filterFunc, filter.Params{}, in[1:], want[1:])
}

func TestEmpty(t *testing.T) {
	filtertest.Check(t, filterFunc, filter.Params{}, "", "+\n+\n")
}
