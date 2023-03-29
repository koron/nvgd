// Package echarts provides chart drawing filter
package echarts

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/render"
	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
)

type Series struct {
	Values []string
}

func (serires *Series) Add(s string) error {
	// FIXME: clever treatment.
	serires.Values = append(serires.Values, s)
	return nil
}

func (serires *Series) AddAll(ss []string) error {
	// FIXME: clever treatment.
	serires.Values = append(serires.Values, ss...)
	return nil
}

type rendererBuilderFunc func([]Series, []charts.GlobalOpts, filter.Params) (render.Renderer, error)

func buildLineRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	if len(serieses) < 2 {
		return nil, fmt.Errorf("less serieses for line chart: %d", len(serieses))
	}
	line := charts.NewLine()
	line.SetGlobalOptions(options...)
	line.SetXAxis(serieses[0].Values[1:])
	for _, s := range serieses[1:] {
		if len(s.Values) == 0 {
			continue
		}
		name := s.Values[0]
		data := make([]opts.LineData, 0, len(s.Values)-1)
		for _, v := range s.Values[1:] {
			dv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			data = append(data, opts.LineData{Value: dv})
		}
		line.AddSeries(name, data)
	}
	return line, nil
}

func buildBarRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	if len(serieses) < 2 {
		return nil, fmt.Errorf("less serieses for bar chart: %d", len(serieses))
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(options...)
	bar.SetXAxis(serieses[0].Values[1:])
	for _, s := range serieses[1:] {
		if len(s.Values) == 0 {
			continue
		}
		name := s.Values[0]
		data := make([]opts.BarData, 0, len(s.Values)-1)
		for _, v := range s.Values[1:] {
			dv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			data = append(data, opts.BarData{Value: dv})
		}
		bar.AddSeries(name, data)
	}
	return bar, nil
}

func buildPieRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	if len(serieses)%2 != 0 {
		return nil, fmt.Errorf("bar chart requires pairs of serieses (name, value): %d", len(serieses))
	}
	if len(serieses[0].Values) != len(serieses[1].Values) {
		return nil, fmt.Errorf("bar chart requires same length of two serieses")
	}
	pie := charts.NewPie()
	pie.SetGlobalOptions(options...)

	for i := 0; i+1 < len(serieses); i += 2 {
		s0, s1 := serieses[i+0], serieses[i+1]
		data := make([]opts.PieData, 0, len(s0.Values)-1)
		for j, n := range s0.Values[1:] {
			v := s1.Values[j+1]
			dv, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, err
			}
			data = append(data, opts.PieData{Name: n, Value: dv})
		}
		name := s0.Values[0]
		if name == "" {
			name = s1.Values[0]
		}
		pie.AddSeries(name, data)
	}

	return pie, nil
}

func buildScatterRenderer(serieses []Series, options []charts.GlobalOpts, p filter.Params) (render.Renderer, error) {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(options...)
	// TODO:
	return scatter, nil
}

var rendererBuilders = map[string]rendererBuilderFunc{
	"line":    buildLineRenderer,
	"bar":     buildBarRenderer,
	"pie":     buildPieRenderer,
	"scatter": buildScatterRenderer,
}

func readSerieses(src *resource.Resource, columnar bool) ([]Series, error) {
	var serieses []Series
	r := csv.NewReader(src)
	var parseColumns func([]string) error = func(ss []string) error {
		if serieses == nil {
			serieses = make([]Series, len(ss))
		}
		for i, s := range ss {
			err := serieses[i].Add(s)
			if err != nil {
				return err
			}
		}
		return nil
	}
	if !columnar {
		parseColumns = func(ss []string) error {
			series := Series{}
			err := series.AddAll(ss)
			if err != nil {
				return err
			}
			serieses = append(serieses, series)
			return nil
		}
	}
	for {
		columns, err := r.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		err = parseColumns(columns)
		if err != nil {
			return nil, err
		}
	}
	return serieses, nil
}

type globalOptionParserFunc func(string) (charts.GlobalOpts, error)

// parseLegendOpts parses a string as opts.Title JSON
//
// See https://pkg.go.dev/github.com/go-echarts/go-echarts/v2@v2.2.5/opts#Title
// for JSON definition
func parseTitleOpts(s string) (charts.GlobalOpts, error) {
	var opt opts.Title
	err := json.Unmarshal([]byte(s), &opt)
	if err != nil {
		return nil, err
	}
	return charts.WithTitleOpts(opt), nil
}

// parseLegendOpts parses a string as opts.Legend JSON
//
// See https://pkg.go.dev/github.com/go-echarts/go-echarts/v2@v2.2.5/opts#Legend
// for JSON definition
func parseLegendOpts(s string) (charts.GlobalOpts, error) {
	var opt opts.Legend
	err := json.Unmarshal([]byte(s), &opt)
	if err != nil {
		return nil, err
	}
	return charts.WithLegendOpts(opt), nil
}

var globalOptionParsers = map[string]globalOptionParserFunc{
	"titleOpts":  parseTitleOpts,
	"legendOpts": parseLegendOpts,
	// FIXME: add parsers for other global option
}

func readGlobalOptions(p filter.Params) ([]charts.GlobalOpts, error) {
	var retval []charts.GlobalOpts
	for k, parser := range globalOptionParsers {
		v, ok := p[k]
		if !ok || v == "" {
			continue
		}
		globalOpts, err := parser(v)
		if err != nil {
			return nil, err
		}
		if globalOpts != nil {
			retval = append(retval, globalOpts)
		}
	}
	return retval, nil
}

func EchartsFilter(src *resource.Resource, p filter.Params) (*resource.Resource, error) {
	var (
		chartType = p.String("t", "line")   // "line", "bar", "pie", "scatter"
		seriesDir = p.String("d", "column") // "column", "row"
	)
	defer src.Close()

	builder, ok := rendererBuilders[strings.ToLower(chartType)]
	if !ok {
		return nil, fmt.Errorf("unknown chart type: %s", chartType)
	}

	_ = seriesDir

	// read "r" as CSV or TSV
	var serieses []Series
	serieses, err := readSerieses(src, strings.ToLower(seriesDir) == "column")
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV/TSV: %w", err)
	}

	// extract global options from "p"
	options, err := readGlobalOptions(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read global options: %w", err)
	}

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
	return src.Wrap(io.NopCloser(outbuf)).PutContentType("text/html"), nil
}

func init() {
	filter.MustRegister("echarts", EchartsFilter)
}
