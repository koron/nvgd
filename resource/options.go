package resource

const (
	// ContentType is for header.
	ContentType = "Content-Type"

	// Filename is for header.
	Filename = "File-Name"
)

// Options holds option properties of Resource.
type Options map[string]interface{}

func (opts Options) clone() Options {
	dist := Options{}
	for k, v := range opts {
		dist[k] = v
	}
	return dist
}

// Bool get value as bool.
func (opts Options) Bool(key string) (value bool, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return false, false
	}
	v, ok := raw.(bool)
	return v, ok
}

// Strings get value as []string.
func (opts Options) String(key string) (value string, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return "", false
	}
	v, ok := raw.(string)
	return v, ok
}

// Strings get value as []string.
func (opts Options) Strings(key string) (value []string, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return nil, false
	}
	v, ok := raw.([]string)
	return v, ok
}
