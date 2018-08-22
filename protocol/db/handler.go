package db

import (
	"errors"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/koron/nvgd/common_const"
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
	query := path(u)
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
	rc, truncated, err := h.execQuery(c, query)
	if err != nil {
		return nil, err
	}
	rs := resource.New(rc)
	rs.Put(common_const.SQLQuery, query)
	if truncated {
		rs.Put(common_const.SQLTruncatedBy, c.maxRows)
	}
	return rs, nil
}

func (h *Handler) openAsset(s string) (*resource.Resource, error) {
	if s == "" {
		s = "index.html"
	}
	if strings.HasPrefix(s, assetPrefix) {
		s = s[len(assetPrefix):]
	}
	f, err := Assets.Open(s)
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
	rows, err := tx.Query(q)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()
	return rows2ltsv(rows, c.maxRows)
}
