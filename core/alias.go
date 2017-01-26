package core

import "strings"

type alias struct {
	from, to string
}

type aliases []alias

// apply aliases for compatibility with koron/night
var defaultAliases = aliases{
	{"files/", "file:///"},
	{"commands/", "command://"},
}

func (a aliases) apply(path string) string {
	for _, n := range a {
		if strings.HasPrefix(path, n.from) {
			return n.to + path[len(n.from):]
		}
	}
	return path
}
