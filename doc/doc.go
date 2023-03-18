// Package doc provides embed.FS for documents
package doc

import "embed"

//go:embed *.md
var content embed.FS
