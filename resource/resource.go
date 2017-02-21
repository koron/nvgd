package resource

import "io"

// Resource packs ReadCloser and its meta info.
type Resource struct {
	io.ReadCloser
}

// New creates a Resource from ReadCloser.
func New(rc io.ReadCloser) *Resource {
	return &Resource{
		ReadCloser: rc,
	}
}

func (r *Resource) Raw() io.ReadCloser {
	return r.ReadCloser
}

func (r *Resource) Wrap(rc io.ReadCloser) *Resource {
	// TODO: copy other properties of r.
	return &Resource{
		ReadCloser: rc,
	}
}

// ReadSeekCloser obtains ReadSeekCloser if it could.
func (r *Resource) ReadSeekCloser() (ReadSeekCloser, bool) {
	x, ok := r.ReadCloser.(ReadSeekCloser)
	return x, ok
}
