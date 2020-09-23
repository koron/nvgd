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
	{"config/", "config://"},
	{"help/", "help://"},
	{"version/", "version://"},
}

func (a aliases) apply(path string) string {
	for _, n := range a {
		if strings.HasPrefix(path, n.from) {
			return n.to + path[len(n.from):]
		}
	}
	return path
}

func (a aliases) mergeMap(m map[string]string) aliases {
	dst := make(aliases, len(a), len(a)+len(m))
	copy(dst[:len(a)], a)
	if len(m) == 0 {
		return dst
	}
	for from, to := range m {
		dst = append(dst, alias{from: from, to: to})
	}
	return dst
}
