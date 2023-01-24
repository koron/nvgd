package db

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/koron/nvgd/internal/commonconst"
	"github.com/koron/nvgd/protocol"
	"github.com/koron/nvgd/resource"
)

const assetPrefix = "assets/"

// Handler is database protocol handler.
type Handler struct {
}

func init() {
	protocol.MustRegister("db", &Handler{})
}

// Open creates a database handler.
func (h *Handler) Open(u *url.URL) (*resource.Resource, error) {
	query := regulatePath(u)
	if query == "" || strings.HasPrefix(query, assetPrefix) {
		return h.openAsset(query)
	}
	c, err := openDB(u)
	if err != nil {
		return nil, err
	}
	if err := h.checkSanity(query); err != nil {
		return nil, err
	}
	st := time.Now()
	rc, truncated, err := h.execQuery(c, query)
	dur := time.Since(st)
	if err != nil {
		return nil, err
	}
	rs := resource.New(rc)
	rs.Put(commonconst.LTSV, true)
	rs.Put(commonconst.SQLQuery, query)
	if truncated {
		rs.Put(commonconst.SQLTruncatedBy, c.maxRows)
	}
	rs.Put(commonconst.SQLExecTime, dur)
	return rs, nil
}

func (h *Handler) openAsset(s string) (*resource.Resource, error) {
	if s == "" {
		s = "index.html"
	}
	s = strings.TrimPrefix(s, assetPrefix)
	f, err := assetsOpen(s)
	if err != nil {
		return nil, err
	}
	rs := resource.New(f).GuessContentType(s)
	return rs, nil
}

var reBadQuery = regexp.MustCompile(`(?i:^\s*(?:insert|update|delete|create|drop|alter|truncate|prepare|execute))`)

func (h *Handler) checkSanity(q string) error {
	// FIXME: too simple, should do more.
	if reBadQuery.MatchString(q) {
		return errors.New("including invalid keywords")
	}
	return nil
}

// execQuery executes a query in a transaction which will be rollbacked.
func (h *Handler) execQuery(c *conn, q string) (io.ReadCloser, bool, error) {
	tx, err := c.db.Begin()
	if err != nil {
		return nil, false, err
	}
	defer tx.Rollback()

	queries := splitQuery(q)
	qlen := len(queries)
	var preQueries []string
	var mainQuery string
	if qlen <= 0 {
		return nil, false, fmt.Errorf("evaluated as empty query: %q", q)
	}
	if qlen == 1 {
		mainQuery = queries[0]
	} else {
		preQueries = queries[:qlen-1]
		mainQuery = queries[qlen]
	}

	// Execute pre-queries.
	for _, query := range preQueries {
		_, err := tx.Exec(query)
		if err != nil {
			return nil, false, fmt.Errorf("pre query %q failed: %w", query, err)
		}
	}

	// `maxRows` limitation will be bypassed when the query has some
	// limitations: "COUNT()" or "LIMIT".
	maxRows := c.maxRows
	if hasLimit(q) {
		maxRows = 0
	}
	// do query.
	rows, err := tx.Query(mainQuery)
	if err != nil {
		return nil, false, fmt.Errorf("main query %q failed as: %w", mainQuery, err)
	}
	defer rows.Close()
	return rows2ltsv(rows, maxRows)
}

var rxSelectCount = regexp.MustCompile(`(?imsU:\bSELECT\b.*\bCOUNT\b.*\(.*\bFROM\b)`)
var rxHasLimit = regexp.MustCompile(`(?imsU:\bLIMIT\b[ ].*\d+)`)

// hasLimit checks a query has LIMIT clause or not.
func hasLimit(q string) bool {
	if rxSelectCount.MatchString(q) {
		return true
	}
	if rxHasLimit.MatchString(q) {
		return true
	}
	return false
}

var rxStripTail = regexp.MustCompile(`\s*;\s*$`)
var rxQuerySplitter = regexp.MustCompile(`\s*;\s*\n\s*`)

// splitQuery splits queries at ';' at end of line.
func splitQuery(s string) []string {
	return rxQuerySplitter.Split(rxStripTail.ReplaceAllString(s, ""), -1)
}
