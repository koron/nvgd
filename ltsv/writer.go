package ltsv

import (
	"bufio"
	"io"
)

// Writer is LTSV writer.
type Writer struct {
	writer *bufio.Writer
	labels []string
}

// NewWriter creates new LTSV writer.
func NewWriter(w io.Writer, labels ...string) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
		labels: labels,
	}
}

// Write writes a LTSV line.
func (w *Writer) Write(values ...string) error {
	for i, l := range w.labels {
		if i != 0 {
			w.writer.WriteRune('\t')
		}
		w.writer.WriteString(l)
		w.writer.WriteRune(':')
		if i < len(values) {
			w.writer.WriteString(values[i])
		}
	}
	w.writer.WriteRune('\n')
	return w.writer.Flush()
}
