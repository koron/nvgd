module github.com/koron/nvgd

go 1.13

require (
	github.com/aws/aws-sdk-go v1.34.29
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/frankban/quicktest v1.7.3 // indirect
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.3.1 // indirect
	github.com/jessevdk/go-assets v0.0.0-20160921144138-4f4301a06e15
	github.com/koron/go-xlsx4db v0.0.3
	github.com/lib/pq v1.8.0
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/olekukonko/tablewriter v0.0.4
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible
	github.com/russross/blackfriday v1.5.2
	github.com/tealeg/xlsx v1.0.5
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

//replace github.com/koron/go-xlsx4db => ../go-xlsx4db
