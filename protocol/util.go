package protocol

import "regexp"

var rxLastComponent = regexp.MustCompile(`[^/]+/?$`)
