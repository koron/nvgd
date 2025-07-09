package pager

import (
	"bytes"
	"fmt"
)

type pageWriter struct {
	num     int
	showNum bool
	written bool
	// w1 write when page matches
	w1 *bytes.Buffer
	// w2 write always if available
	w2 *bytes.Buffer
}

func (pw *pageWriter) write(b []byte) (int, error) {
	if pw.w1 != nil {
		pw.w1.Write(b)
	}
	if pw.w2 != nil {
		pw.w2.Write(b)
	}
	return len(b), nil
}

func (pw *pageWriter) Write(b []byte) (int, error) {
	if !pw.written {
		pw.written = true
		if pw.showNum {
			fmt.Fprintf(pw, "(page %d)\n", pw.num)
		}
	}
	return pw.write(b)
}
