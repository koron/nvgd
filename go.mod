module github.com/koron/nvgd

go 1.21

require (
	github.com/aws/aws-sdk-go v1.55.5
	github.com/go-echarts/go-echarts/v2 v2.4.1
	github.com/go-redis/redis/v7 v7.4.1
	github.com/go-sql-driver/mysql v1.8.1
	github.com/google/go-cmp v0.6.0
	github.com/koron-go/ringbuf v1.0.0
	github.com/koron/go-xlsx4db v0.0.3
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.23
	github.com/noborus/trdsql v0.11.1
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pierrec/lz4/v4 v4.1.21
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/tealeg/xlsx v1.0.5
	gopkg.in/yaml.v2 v2.4.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/itchyny/gojq v0.12.12 // indirect
	github.com/itchyny/timefmt-go v0.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jwalton/gchalk v1.3.0 // indirect
	github.com/jwalton/go-supportscolor v1.1.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/multiprocessio/go-sqlite3-stdlib v0.0.0-20220822170115-9f6825a1cd25 // indirect
	github.com/noborus/guesswidth v0.3.1 // indirect
	github.com/noborus/tbln v0.0.2 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/term v0.24.0 // indirect
	gonum.org/v1/gonum v0.12.0 // indirect
)

replace (
	//github.com/koron/go-xlsx4db => ../go-xlsx4db
	github.com/noborus/trdsql v0.11.1 => ./_replace/trdsql@v0.11.1
	github.com/russross/blackfriday/v2 v2.1.0 => github.com/koron/blackfriday/v2 v2.1.0-fix.2
)
