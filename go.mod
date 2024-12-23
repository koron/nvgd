module github.com/koron/nvgd

go 1.22.10

require (
	github.com/alicebob/miniredis/v2 v2.34.0
	github.com/aws/aws-sdk-go v1.55.5
	github.com/go-echarts/go-echarts/v2 v2.4.1
	github.com/go-redis/redis/v7 v7.4.1
	github.com/go-sql-driver/mysql v1.8.1
	github.com/google/go-cmp v0.6.0
	github.com/koron-go/ringbuf v1.0.0
	github.com/koron/go-xlsx4db v0.0.3
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/noborus/trdsql v1.1.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pierrec/lz4/v4 v4.1.22
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/tealeg/xlsx v1.0.5
	gopkg.in/yaml.v2 v2.4.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/alicebob/gopher-json v0.0.0-20230218143504-906a9b012302 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/goccy/go-yaml v1.11.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/iancoleman/orderedmap v0.3.0 // indirect
	github.com/itchyny/gojq v0.12.16 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jwalton/gchalk v1.3.0 // indirect
	github.com/jwalton/go-supportscolor v1.2.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/multiprocessio/go-sqlite3-stdlib v0.0.0-20220822170115-9f6825a1cd25 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/noborus/guesswidth v0.4.0 // indirect
	github.com/noborus/sqlss v0.1.0 // indirect
	github.com/noborus/tbln v0.0.2 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/yuin/gopher-lua v1.1.1 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/term v0.27.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	gonum.org/v1/gonum v0.14.0 // indirect
	modernc.org/gc/v3 v3.0.0-20240304020402-f0dba7c97c2b // indirect
	modernc.org/libc v1.55.3 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/sqlite v1.31.1 // indirect
	modernc.org/strutil v1.2.0 // indirect
	modernc.org/token v1.1.0 // indirect
)

replace github.com/russross/blackfriday/v2 v2.1.0 => github.com/koron/blackfriday/v2 v2.1.0-fix.2
