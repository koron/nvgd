// Package fileentry provides file entry related operations.
package fileentry

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/koron/nvgd/internal/ltsv"
)

type Entry struct {
	Name       string
	Type       string
	Size       int64
	ModifiedAt time.Time
	Link       string
	Download   string

	useUnixTime bool
}

func (e Entry) IsDir() bool {
	return e.Type == "dir" || e.Type == "prefix"
}

func parseTime(s string) (time.Time, bool, error) {
	sec, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return time.Unix(sec, 0), true, nil
	}
	ti, err := time.Parse(time.RFC1123, s)
	return ti, false, err
}

func ParseLTSV(s *ltsv.Set) (Entry, error) {
	entry := Entry{
		Name:     s.GetFirst("name"),
		Type:     s.GetFirst("type"),
		Link:     s.GetFirst("link"),
		Download: s.GetFirst("download"),
	}
	switch entry.Type {
	case "dir", "prefix", "file", "object":
		// Valid type
	default:
		return Entry{}, fmt.Errorf("unknown entry type: %s", entry.Type)
	}

	nsize, _ := strconv.ParseInt(s.GetFirst("size"), 10, 64)
	entry.Size = nsize

	if v := s.GetFirst("modified_at"); v != "" {
		modtime, useUnixTime, err := parseTime(v)
		if err != nil {
			return Entry{}, err
		}
		entry.ModifiedAt = modtime
		entry.useUnixTime = useUnixTime
	}

	return entry, nil
}

type LTSVWriter struct {
	w *ltsv.Writer

	UseUnixTime bool
}

func NewLTSVWriter(w io.Writer) *LTSVWriter {
	ltsvW := ltsv.NewWriter(w, "name", "type", "size", "modified_at", "link", "download")
	return &LTSVWriter{w: ltsvW}
}

func (w *LTSVWriter) timeStr(ti time.Time, useUnixTime bool) string {
	if ti.IsZero() {
		return ""
	}
	if w.UseUnixTime || useUnixTime {
		return strconv.FormatInt(ti.Unix(), 10)
	}
	return ti.Format(time.RFC1123)
}

func (w *LTSVWriter) Write(e Entry) error {
	return w.w.Write(
		e.Name,
		e.Type,
		strconv.FormatInt(e.Size, 10),
		w.timeStr(e.ModifiedAt, e.useUnixTime),
		e.Link,
		e.Download)
}
