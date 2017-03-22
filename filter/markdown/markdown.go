package markdown

import (
	"bytes"
	"io/ioutil"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
	"github.com/russross/blackfriday"
)

func init() {
	filter.MustRegister("markdown", filterMarkdown)
}

func filterMarkdown(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// TODO:
	b1, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b2 := blackfriday.MarkdownCommon(b1)
	r2 := ioutil.NopCloser(bytes.NewReader(b2))
	return r.Wrap(r2), nil
}
