package core

import (
	"testing"

	"github.com/koron/nvgd/config"
	"github.com/koron/nvgd/internal/assert"
)

func TestGetRawDeterministic(t *testing.T) {
	f := &Filters{
		descs: config.FiltersMap{
			"file:///var/":     {"tail"},
			"file:///var/log/": {"head"},
		},
	}

	// Run multiple times to ensure determinism.
	for i := 0; i < 100; i++ {
		got, found := f.getRaw("file:///var/log/messages")
		if !found {
			t.Fatalf("iteration %d: expected to find a match", i)
		}
		// Longest prefix "file:///var/log/" should match.
		if len(got) != 1 || got[0] != "head" {
			t.Fatalf("iteration %d: expected 'head' (from longest prefix), got %v", i, got)
		}
	}
}

func TestGetRawNoMatch(t *testing.T) {
	f := &Filters{
		descs: config.FiltersMap{
			"file:///var/log/": {"tail"},
		},
	}
	_, found := f.getRaw("command://uptime")
	assert.Equal(t, false, found, "should not match for unrelated path")
}
