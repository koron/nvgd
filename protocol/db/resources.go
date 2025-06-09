package db

import (
	"embed"
	"io/fs"
	"path"

	"github.com/koron/nvgd/internal/devfs"
)

//go:embed assets
var embedFS embed.FS

var assetsFS = devfs.New(embedFS, "protocol/db", "")

func assetsOpen(name string) (fs.File, error) {
	return assetsFS.Open(path.Join("assets", name))
}
