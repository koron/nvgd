package db

import (
	"embed"
	"io/fs"
	"path"
)

//go:embed assets
var assetsFS embed.FS

func assetsOpen(name string) (fs.File, error) {
	return assetsFS.Open(path.Join("assets", name))
}
