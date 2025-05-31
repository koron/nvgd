package plugins

import (
	_ "github.com/koron/nvgd/protocol/aws"
	_ "github.com/koron/nvgd/protocol/command"
	_ "github.com/koron/nvgd/protocol/configp"
	_ "github.com/koron/nvgd/protocol/db"
	_ "github.com/koron/nvgd/protocol/duckdb"
	_ "github.com/koron/nvgd/protocol/echarts"
	_ "github.com/koron/nvgd/protocol/examples"
	_ "github.com/koron/nvgd/protocol/file"
	_ "github.com/koron/nvgd/protocol/help"
	_ "github.com/koron/nvgd/protocol/redis"
	_ "github.com/koron/nvgd/protocol/trdsql"
	_ "github.com/koron/nvgd/protocol/version"
)
