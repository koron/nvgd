package hash

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestHashFilter(t *testing.T) {
	// encoding variations
	filtertest.Check(t, newHash,
		filter.Params{},
		"",
		"d41d8cd98f00b204e9800998ecf8427e")
	filtertest.Check(t, newHash,
		filter.Params{ "encoding": "base64", },
		"",
		"1B2M2Y8AsgTpgAmY7PhCfg==")
	filtertest.Check(t, newHash,
		filter.Params{ "encoding": "binary", },
		"",
		"\xd4\x1d\x8cُ\x00\xb2\x04\xe9\x80\t\x98\xec\xf8B~")

	// algorithm variations
	filtertest.Check(t, newHash,
		filter.Params{ "algorithm": "sha1"},
		"",
		"da39a3ee5e6b4b0d3255bfef95601890afd80709")
	filtertest.Check(t, newHash,
		filter.Params{ "algorithm": "sha1"},
		"",
		"da39a3ee5e6b4b0d3255bfef95601890afd80709")
	filtertest.Check(t, newHash,
		filter.Params{ "algorithm": "sha256"},
		"",
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	filtertest.Check(t, newHash,
		filter.Params{ "algorithm": "sha512"},
		"",
		"cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e")
}
