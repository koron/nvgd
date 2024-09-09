// Package commonconst provides constant values for NVGD.
package commonconst

const (
	// Small indicates that the content is shortened. Default filters are not
	// applied to content with this tag.
	Small = "x-small-source"

	// ParsedKeys is a tag that indicates the name of the query parameter as
	// interpreted by the protocol. The value of this tag is a string array,
	// and the query parameter names contained therein are excluded before the
	// filter is applied. This means that the query parameters are consumed by
	// the protocol rather than the filter.
	ParsedKeys = "x-parsed-keys"
)
