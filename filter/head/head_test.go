package head

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestHead(t *testing.T) {
	filtertest.Check(t, newHead,
		filter.Params{},
		"0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n",
		"0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n")
	filtertest.Check(t, newHead,
		filter.Params{
			"start": "3",
		},
		"0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n",
		"3\n4\n5\n6\n7\n8\n9\n")
	filtertest.Check(t, newHead,
		filter.Params{
			"start": "3",
			"limit": "5",
		},
		"0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n",
		"3\n4\n5\n6\n7\n")
}
