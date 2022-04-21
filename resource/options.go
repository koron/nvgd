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

// Bool gets a value as bool.
func (opts Options) Bool(key string) (value bool, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return false, false
	}
	v, ok := raw.(bool)
	return v, ok
}

// String gets a value as string.
func (opts Options) String(key string) (value string, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return "", false
	}
	v, ok := raw.(string)
	return v, ok
}

// Strings gets a value as []string.
func (opts Options) Strings(key string) (value []string, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return nil, false
	}
	v, ok := raw.([]string)
	return v, ok
}

// Int gets a value as int.
func (opts Options) Int(key string) (value int, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return 0, false
	}
	v, ok := raw.(int)
	return v, ok
}
