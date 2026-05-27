package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"

	"github.com/koron/nvgd/internal/ltsv"
)

// NullReplacement replaces null value in LTSV.
var NullReplacement = "(null)"

func rows2ltsv(rows *sql.Rows, maxRows int) (io.ReadCloser, bool, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, false, err
	}
	var (
		buf = &bytes.Buffer{}
		w   = ltsv.NewWriter(buf, cols...)
		n   = len(cols)
	)

	vals := make([]any, n)
	for i := range vals {
		vals[i] = new(any)
	}
	strs := make([]string, n)

	nrow := 0
	truncated := false
	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return nil, false, err
		}
		for i, v := range vals {
			p := v.(*any)
			switch val := (*p).(type) {
			case nil:
				strs[i] = NullReplacement
			case []byte:
				strs[i] = string(val)
			default:
				strs[i] = fmt.Sprint(val)
			}
		}
		w.Write(strs...)
		nrow++
		if maxRows > 0 && nrow >= maxRows {
			truncated = rows.Next()
			if err := rows.Err(); err != nil {
				return nil, false, err
			}
			break
		}
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	return io.NopCloser(buf), truncated, nil
}
