// Package plugins load default protocols and filters for nvgd.
package plugins

import (
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
	_ "github.com/koron/nvgd/protocol/help"
	_ "github.com/koron/nvgd/protocol/redis"
	_ "github.com/koron/nvgd/protocol/version"
)
