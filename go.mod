module github.com/koron/nvgd

go 1.19

require (
	github.com/aws/aws-sdk-go v1.44.209
	github.com/go-redis/redis/v7 v7.4.1
	github.com/go-sql-driver/mysql v1.7.0
	github.com/google/go-cmp v0.5.9
	github.com/koron/go-xlsx4db v0.0.3
	github.com/lib/pq v1.10.7
	github.com/olekukonko/tablewriter v0.0.5
	github.com/pierrec/lz4/v4 v4.1.17
	github.com/russross/blackfriday v1.6.0
	github.com/tealeg/xlsx v1.0.5
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	golang.org/x/crypto v0.6.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
)

//replace github.com/koron/go-xlsx4db => ../go-xlsx4db
