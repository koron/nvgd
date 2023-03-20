// Package echarts provides chart drawing filter
package echarts

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
)

type Series struct{}

type rendererBuilder func(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error)

func buildLineRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	line := charts.NewLine()
	line.SetGlobalOptions(options...)
	// TODO:
	return line, nil
}

func buildBarRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	bar := charts.NewBar()
	bar.SetGlobalOptions(options...)
	// TODO:
	return bar, nil
}

func buildPieRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(options...)
	// TODO:
	return pie, nil
}

func buildScatterRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(options...)
	// TODO:
	return scatter, nil
}

var rendererBuilders = map[string]rendererBuilder{
	"line":    buildLineRenderer,
	"bar":     buildBarRenderer,
	"pie":     buildPieRenderer,
	"scatter": buildScatterRenderer,
}

func EchartsFilter(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	var (
		chartType = p.String("t", "line")   // "line", "bar", "pie", "scatter"
		seriesDir = p.String("d", "column") // "column", "row"
	)

	builder, ok := rendererBuilders[strings.ToLower(chartType)]
	if !ok {
		return nil, fmt.Errorf("unknown chart type: %s", chartType)
	}

	_ = seriesDir

	// TODO: read "r" as CSV or TSV
	var serieses []Series

	// FIXME: extract global options from "p"
	var options []charts.GlobalOpts = []charts.GlobalOpts{}

	renderer, err := builder(serieses, options, p)
	if err != nil {
		return nil, fmt.Errorf("failed to build chart renderer: %w", err)
	}

	// render a chart
	outbuf := &bytes.Buffer{}
	err = renderer.Render(outbuf)
	if err != nil {
		return nil, fmt.Errorf("failed to render chart: %w", err)
	}
	return r.Wrap(io.NopCloser(outbuf)).PutContentType("text/html"), nil
}

func init() {
	filter.MustRegister("echarts", EchartsFilter)
}
