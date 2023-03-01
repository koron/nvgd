// Package trdsql provides https://github.com/noborus/trdsql filter
package trdsql

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/koron/nvgd/filter"
	"github.com/koron/nvgd/resource"
	"github.com/noborus/trdsql"
)

var execTimeout time.Duration = 30 * time.Second

func trdsqlFilter(r *resource.Resource, p filter.Params) (*resource.Resource, error) {
	// process parameters
	var (
		inName = p.String("iname", "t")

		inFormat    = parseInFormat(p.String("ifmt", ""), trdsql.CSV)
		inDelim     = p.String("id", ",")
		inHeader    = p.Bool("ih", false)
		inSkip      = p.Int("is", 0)
		inPreRead   = p.Int("ir", 1)
		inLimitRead = p.Bool("il", false)
		inJQ        = p.String("ijq", "")

		inNULL, inNeedNULL = p["inull"]

		outFormat    = parseOutFormat(p.String("ofmt", ""), trdsql.CSV)
		outDelim     = p.String("od", ",")
		outQuote     = p.String("oq", "\"")
		outAllQuotes = p.Bool("oaq", false)
		outUseCRLF   = p.Bool("ocrlf", false)
		outHeader    = p.Bool("oh", false)
		outNoWrap    = p.Bool("onowrap", false)

		outNULL, outNeedNULL = p["onull"]

		query = p["q"]
	)

	// setup importer
	stdin := r
	importer, err := trdsql.NewBufferImporter(inName, stdin,
		trdsql.InFormat(inFormat),
		trdsql.InDelimiter(inDelim),
		trdsql.InHeader(inHeader),
		trdsql.InSkip(inSkip),
		trdsql.InPreRead(inPreRead),
		trdsql.InLimitRead(inLimitRead),
		trdsql.InJQ(inJQ),
		trdsql.InNeedNULL(inNeedNULL),
		trdsql.InNULL(inNULL),
	)
	if err != nil {
		return nil, fmt.Errorf("trdsql: failed to create importer: %w", err)
	}

	// setup exporter
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	exporter := trdsql.NewExporter(trdsql.NewWriter(
		trdsql.OutFormat(outFormat),
		trdsql.OutDelimiter(outDelim),
		trdsql.OutQuote(outQuote),
		trdsql.OutAllQuotes(outAllQuotes),
		trdsql.OutUseCRLF(outUseCRLF),
		trdsql.OutHeader(outHeader),
		trdsql.OutNoWrap(outNoWrap),
		trdsql.OutNeedNULL(outNeedNULL),
		trdsql.OutNULL(outNULL),
		trdsql.OutStream(stdout),
		trdsql.ErrStream(stderr),
	))

	// execute a query
	trd := trdsql.NewTRDSQL(importer, exporter)
	ctx, cancel := context.WithTimeout(context.Background(), execTimeout)
	defer cancel()
	err = trd.ExecContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("trdsql error: %w: %s", err, stderr.String())
	}
	return r.Wrap(io.NopCloser(stdout)), nil
}

func parseInFormat(s string, defaultFormat trdsql.Format) trdsql.Format {
	switch strings.ToUpper(s) {
	case "CSV":
		return trdsql.CSV
	case "LTSV":
		return trdsql.LTSV
	case "JSON":
		return trdsql.JSON
	case "TBLN":
		return trdsql.TBLN
	case "TSV":
		return trdsql.TSV
	case "PSV":
		return trdsql.PSV
	default:
		return defaultFormat
	}
}

func parseOutFormat(s string, defaultFormat trdsql.Format) trdsql.Format {
	switch strings.ToUpper(s) {
	case "CSV":
		return trdsql.CSV
	case "LTSV":
		return trdsql.LTSV
	case "JSON":
		return trdsql.JSON
	case "TBLN":
		return trdsql.TBLN
	case "RAW":
		return trdsql.RAW
	case "MD":
		return trdsql.MD
	case "AT":
		return trdsql.AT
	case "JSONL":
		return trdsql.JSONL
	default:
		return defaultFormat
	}
}

func init() {
	filter.MustRegister("trdsql", trdsqlFilter)
}
