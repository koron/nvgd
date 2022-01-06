package main

import (
	_ "embed"

	_ "github.com/koron/nvgd/filter/count"
	_ "github.com/koron/nvgd/filter/htmltable"
	_ "github.com/koron/nvgd/filter/indexhtml"
	_ "github.com/koron/nvgd/filter/jsonarray"
	_ "github.com/koron/nvgd/filter/markdown"
	_ "github.com/koron/nvgd/filter/tail"
	_ "github.com/koron/nvgd/filter/texttable"
	_ "github.com/koron/nvgd/protocol/aws"
	_ "github.com/koron/nvgd/protocol/configp"
	_ "github.com/koron/nvgd/protocol/db"
	"github.com/koron/nvgd/protocol/help"
	_ "github.com/koron/nvgd/protocol/redis"
	_ "github.com/koron/nvgd/protocol/version"
)

//go:embed README.md
var readme string

func init() {
	help.Text = readme
}
