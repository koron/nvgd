package ltsv

import (
	"bufio"
	"bytes"
	"io"
	"strings"
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
			w.writer.WriteString(escape(values[i]))
		}
	}
	w.writer.WriteRune('\n')
	return w.writer.Flush()
}

func escape(s string) string {
	if !strings.ContainsAny(s, "\\\t\n\r") {
		return s
	}
	bb := &bytes.Buffer{}
	for _, r := range s {
		switch r {
		case '\\':
			bb.WriteString(`\\`)
		case '\t':
			bb.WriteString(`\t`)
		case '\n':
			bb.WriteString(`\n`)
		case '\r':
			bb.WriteString(`\r`)
		default:
			bb.WriteRune(r)
		}
	}
	return bb.String()
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
