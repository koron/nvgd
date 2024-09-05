package ltsv

import (
	"bufio"
	"bytes"
	"io"

	"github.com/koron/nvgd/internal/filterbase"
)

// Reader reads LTSV values.
type Reader struct {
	rd *bufio.Reader
}

// NewReader creates a new LTSV reader.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		rd: bufio.NewReaderSize(r, filterbase.Config.MaxLineLen),
	}
}

func (r *Reader) readLine() ([]byte, error) {
	// FIXME: Avoids loading a file with very long line and allocating
	// excessive memory.  It would be better to use a larger buffer from the
	// beginning, nor never extending .
	d, err := r.rd.ReadSlice('\n')
	if err == nil || (err == io.EOF && len(d) > 0) {
		return d, nil
	} else if err != bufio.ErrBufferFull {
		return nil, err
	}
	return nil, filterbase.ErrMaxLineExceeded
}

// Read read a LTSV value.
func (r *Reader) Read() (*Set, error) {
	d, err := r.readLine()
	if err != nil {
		return nil, err
	}
	d = bytes.TrimLeft(d, " \n\r\t")
	d = bytes.TrimRight(d, "\n\r\t")
	s := &Set{
		Index: make(map[string][]int),
	}
	for _, raw := range bytes.Split(d, []byte("\t")) {
		kv := bytes.SplitN(raw, []byte(":"), 2)
		if len(kv) != 2 {
			continue
		}
		s.Put(string(kv[0]), string(kv[1]))
	}
	return s, nil
}

// Set is a set of LTSV values in a line.
type Set struct {
	Properties []Property
	Index      map[string][]int
}

// Put puts a property to the set.
func (s *Set) Put(label, value string) {
	n := len(s.Properties)
	s.Properties = append(s.Properties, Property{Label: label, Value: value})
	s.Index[label] = append(s.Index[label], n)
}

// Get gets values for the label.
func (s *Set) Get(label string) []string {
	indexes, ok := s.Index[label]
	if !ok {
		return nil
	}
	list := make([]string, len(indexes))
	for i, n := range indexes {
		list[i] = s.Properties[n].Value
	}
	return list
}

// GetFirst gets a first value for the label.
func (s *Set) GetFirst(label string) string {
	indexes, ok := s.Index[label]
	if !ok || len(indexes) == 0 {
		return ""
	}
	return s.Properties[indexes[0]].Value
}

// Empty checks set is empty or not.
func (s *Set) Empty() bool {
	return len(s.Properties) == 0
}

// Property is a pair of label and value.
type Property struct {
	Label string
	Value string
}
