package echarts

import (
	"testing"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/internal/filtertest"
)

func TestPieEmptySeries(t *testing.T) {
	filtertest.Fail(t, EchartsFilter,
		filter.Params{"t": "pie"},
		"",
		"failed to build chart renderer: less serieses for pie chart: 0")
}

func TestPieSingleSeries(t *testing.T) {
	filtertest.Fail(t, EchartsFilter,
		filter.Params{"t": "pie"},
		"a",
		"failed to build chart renderer: less serieses for pie chart: 1")
}

func TestPieOddSeries(t *testing.T) {
	filtertest.Fail(t, EchartsFilter,
		filter.Params{"t": "pie"},
		"a,b,c",
		"failed to build chart renderer: pie chart requires pairs of serieses (name, value): 3")
}

func TestPieSeriesLengthMismatch(t *testing.T) {
	s0 := Series{Values: []string{"name", "a", "b"}}
	s1 := Series{Values: []string{"val", "1"}}
	_, err := buildPieRenderer([]Series{s0, s1}, nil, nil)
	if err == nil {
		t.Fatal("expected error for mismatched series lengths")
	}
}

func TestPieSecondPairLengthMismatch(t *testing.T) {
	s0 := Series{Values: []string{"x", "a"}}
	s1 := Series{Values: []string{"y", "1"}}
	s2 := Series{Values: []string{"name", "p", "q", "r"}}
	s3 := Series{Values: []string{"val", "10"}}
	_, err := buildPieRenderer([]Series{s0, s1, s2, s3}, nil, nil)
	if err == nil {
		t.Fatal("expected error for mismatched series lengths in second pair")
	}
}
