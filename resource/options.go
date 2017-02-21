package resource

// Options holds option properties of Resource.
type Options map[string]interface{}

func (opts Options) clone() Options {
	dist := Options{}
	for k, v := range opts {
		dist[k] = v
	}
	return dist
}

// Bool get boolean
func (opts Options) Bool(key string) (value bool, ok bool) {
	raw, ok := opts[key]
	if !ok {
		return false, false
	}
	v, ok := raw.(bool)
	return v, ok
}
