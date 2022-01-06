module github.com/koron/nvgd

go 1.16

require (
	github.com/aws/aws-sdk-go v1.42.28
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/frankban/quicktest v1.7.3 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/koron/go-xlsx4db v0.0.3
	github.com/lib/pq v1.10.4
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/olekukonko/tablewriter v0.0.5
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible
	github.com/russross/blackfriday v1.6.0
	github.com/tealeg/xlsx v1.0.5
	gopkg.in/yaml.v2 v2.4.0
)

//replace github.com/koron/go-xlsx4db => ../go-xlsx4db
