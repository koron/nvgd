package db

import (
	"database/sql"
	"io"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func openMemory(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open :memory: db: %s", err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestRows2LTSV(t *testing.T) {
	db := openMemory(t)

	_, err := db.Exec(`CREATE TABLE t (
		id INTEGER PRIMARY KEY,
		name TEXT,
		blob_data BLOB,
		nullable TEXT,
		ival INTEGER,
		rval REAL
	)`)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}

	_, err = db.Exec(`INSERT INTO t VALUES
		(1, 'Alice', x'deadbeef', NULL, 42, 3.14),
		(2, 'Bob',   x'00ff',     'hello', 0, 0.0),
		(3, '',      x'',         '',     -1, -1.5)`)
	if err != nil {
		t.Fatalf("failed to insert: %s", err)
	}

	rows, err := db.Query(`SELECT id, name, blob_data, nullable, ival, rval FROM t ORDER BY id`)
	if err != nil {
		t.Fatalf("failed to query: %s", err)
	}
	defer rows.Close()

	rc, truncated, err := rows2ltsv(rows, 0)
	if err != nil {
		t.Fatalf("rows2ltsv failed: %s", err)
	}
	if truncated {
		t.Error("unexpected truncated")
	}
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read output: %s", err)
	}
	got := string(b)

	// Check all three rows are present.
	if !strings.Contains(got, "Alice") {
		t.Errorf("output should contain Alice:\n%s", got)
	}
	if !strings.Contains(got, "Bob") {
		t.Errorf("output should contain Bob:\n%s", got)
	}

	// Blob data should be preserved as raw bytes (not corrupted).
	// x'deadbeef' → []byte{0xde, 0xad, 0xbe, 0xef}
	if !strings.Contains(got, "\xde\xad\xbe\xef") {
		t.Errorf("blob data should be preserved as raw bytes:\n% x", []byte(got))
	}
	if !strings.Contains(got, "\x00\xff") {
		t.Errorf("blob with null byte should be preserved:\n% x", []byte(got))
	}

	// NULL should be replaced.
	if !strings.Contains(got, "(null)") {
		t.Errorf("NULL value should be replaced with (null):\n%s", got)
	}

	// Empty string should be preserved (row 3 has empty name → "name:\t").
	if !strings.Contains(got, "\tname:\t") {
		t.Errorf("empty string should be preserved as 'name:' directly followed by tab:\n% x", []byte(got))
	}
}

func TestRows2LTSVTruncated(t *testing.T) {
	db := openMemory(t)

	_, err := db.Exec(`CREATE TABLE t (id INTEGER PRIMARY KEY, val TEXT)`)
	if err != nil {
		t.Fatalf("failed to create table: %s", err)
	}
	for i := 0; i < 10; i++ {
		_, err := db.Exec(`INSERT INTO t VALUES (?, ?)`, i, i)
		if err != nil {
			t.Fatalf("failed to insert: %s", err)
		}
	}
	rows, err := db.Query(`SELECT id, val FROM t ORDER BY id`)
	if err != nil {
		t.Fatalf("failed to query: %s", err)
	}
	defer rows.Close()

	rc, truncated, err := rows2ltsv(rows, 3)
	if err != nil {
		t.Fatalf("rows2ltsv failed: %s", err)
	}
	b, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("failed to read output: %s", err)
	}
	if !truncated {
		t.Errorf("expected truncated=true for maxRows=3")
	}
	got := string(b)
	if strings.Count(got, "\n") != 3 { // 3 data rows (LTSV has no separate header)
		t.Errorf("expected 3 lines for maxRows=3, got %d lines:\n%s",
			strings.Count(got, "\n"), got)
	}
}
