package db

import (
	"bytes"
	"database/sql"
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
		vals[i] = new(sql.NullString)
	}
	strs := make([]string, n)

	nrow := 0
	truncated := false
	for rows.Next() {
		if err := rows.Scan(vals...); err != nil {
			return nil, false, err
		}
		for i, v := range vals {
			ns := v.(*sql.NullString)
			if ns.Valid {
				strs[i] = ns.String
			} else {
				strs[i] = NullReplacement
			}
		}
		w.Write(strs...)
		nrow++
		if maxRows > 0 && nrow >= maxRows {
			truncated = rows.Next()
			break
		}
	}
	return io.NopCloser(buf), truncated, nil
}
