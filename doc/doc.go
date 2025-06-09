// Package doc provides embed.FS for documents
package doc

import (
	"embed"
	"io/fs"

	"github.com/koron/nvgd/internal/devfs"
)

//go:embed *.md
var embedFS embed.FS

// FS provides markdown content in doc/ directory
var FS fs.FS = devfs.New(embedFS, "doc", "*.md")
