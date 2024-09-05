package grep

import "bytes"

// LineFilter is a line filter.
type LineFilter func([]byte) []byte

// Apply applies a filter to a line bytes.
func (lf LineFilter) Apply(b []byte) []byte {
	if lf == nil {
		return b
	}
	return lf(b)
}

// Chain createa a new LineFilters which chains two LineFilters.
func (lf LineFilter) Chain(second LineFilter) LineFilter {
	if lf == nil {
		return second
	}
	if second == nil {
		return lf
	}
	return func(b []byte) []byte {
		return second.Apply(lf.Apply(b))
	}
}

// TrimEOL trims EOL bytes.
var TrimEOL LineFilter = func(b []byte) []byte {
	l := len(b)
	if l >= 1 && b[l-1] == '\n' {
		c := 1
		if l >= 2 && b[len(b)-2] == '\r' {
			c++
		}
		return b[:len(b)-c]
	}
	return b
}

// NewCutLF creates a cut line filter with sep and position.
func NewCutLF(sep []byte, n int) LineFilter {
	m := n + 2
	return func(b []byte) []byte {
		fields := bytes.SplitN(b, sep, m)
		if n >= len(fields) {
			return b
		}
		return fields[n]
	}
}
