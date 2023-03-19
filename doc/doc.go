// Package doc provides embed.FS for documents
package doc

import "embed"

// FS provides markdown content in doc/ directory
//go:embed *.md
var FS embed.FS
