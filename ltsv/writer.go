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

func Write(w io.StringWriter, props []Property) error {
	for i, p := range props {
		if i != 0 {
			_, err := w.WriteString("\t")
			if err != nil {
				return err
			}
		}
		_, err := w.WriteString(p.Label)
		if err != nil {
			return err
		}
		_, err = w.WriteString(":")
		if err != nil {
			return err
		}
		_, err = w.WriteString(p.Value)
		if err != nil {
			return err
		}
	}
	_, err := w.WriteString("\n")
	return err
}
