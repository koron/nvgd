package db

import (
	"testing"

	"github.com/koron/nvgd/internal/assert"
)

func TestBadQuery(t *testing.T) {
	f := func(s string, expected bool) {
		got := checkSQLSanity(s)
		if expected && got == nil {
			t.Errorf("checkSQLSanity(%q) should be an error, got nil", s)
		}
		if !expected && got != nil {
			t.Errorf("checkSQLSanity(%q) should be nil, got %v", s, got)
		}
	}
	// basic keywords
	f("INSERT", true)
	f("UPDATE", true)
	f("DELETE", true)
	f("CREATE", true)
	f("DROP", true)
	f("ALTER", true)
	f("TRUNCATE", true)
	f("EXECUTE", true)
	f("PREPARE", true)
	f("REPLACE", true)
	// variations
	f("insert", true)
	f("UpDaTe", true)
	f(" DELETE", true)
	f(" CrEaTe", true)
	f(" DROP", true)
	f(" ALTer", true)
	// multi-statement injections
	f("SELECT 1; DROP TABLE t", true)
	// comment bypass — block comments
	f("INS/**/ERT INTO t", true)
	f("/**/INSERT INTO t", true)
	f("INSERT/**/INTO t VALUES(1)", true)
	f("SEL/**/ECT 1", false) // SEL/**/ECT is not a bad keyword
	// comment bypass — line comments
	f("-- comment\nINSERT INTO t", true)
	f("SELECT 1; -- comment\nDROP TABLE t", true)
	f("SELECT 1;\n-- comment\nDROP TABLE t", true)
	// inhibits
	f("SELECT * FROM USERS", false)
	f("SELECT * FROM UPDATES", false)
	f("SELECT * FROM INSERTS", false)
	f("SELECT * FROM DROP_TABLE", false)
	// subquery with bad keyword in string literal — not blocked (acceptable)
	f("SELECT * FROM t WHERE name = 'INSERT'", false)
}

func TestHasLimit(t *testing.T) {
	for i, tc := range []struct {
		query string
		want  bool
	}{
		{``, false},
		{`SELECT gid, COUNT(*) FROM user`, true}, // SELECT with COUNT
		{`select gid, COUNT(*) from user`, true}, // case insensitive
		{`SELECT * FROM user LIMIT 10000`, true},   // LIMIT with number
		{`LIMIT 10000`, true},                      // only LIMIT with number
		{`SELECT * FROM user LIMIT ?`, true},       // LIMIT with placeholder
		{`SELECT * FROM user LIMIT @limit`, true},  // LIMIT with named parameter
		{`SELECT * FROM user LIMIT :limit`, true},  // LIMIT with named parameter
		{`SELECT * FROM user LIMIT abc`, true},     // any token after LIMIT is considered a limit clause
		{"LIMIT\n10000", true},                     // LIMIT with a new line
		{"LIMIT \n10000", true},                    // LIMIT with a new line
		{"LIMIT\n 10000", true},                    // LIMIT with a new line
		{"LIMIT\n\n10000", true},                   // LIMIT with new lines
		{"LIMIT\n\n\n10000", true},                 // LIMIT with new lines
	} {
		got := hasLimit(tc.query)
		if got != tc.want {
			t.Errorf("unexpected hasLimit return: want=%t got=%t:#%d: %q", tc.want, got, i, tc.query)
		}
	}
}

func TestSplitQuery(t *testing.T) {
	for i, tc := range []struct {
		in   string
		want []string
	}{
		{"SELECT * FROM ACCOUNT", []string{"SELECT * FROM ACCOUNT"}},
		{"SELECT * FROM ACCOUNT;", []string{"SELECT * FROM ACCOUNT"}},
		{
			"SET FOOBAR=123 ; \n SELECT * FROM ACCOUNT ;",
			[]string{
				"SET FOOBAR=123",
				"SELECT * FROM ACCOUNT",
			},
		},
	} {
		got := splitQuery(tc.in)
		assert.Equal(t, tc.want, got, "unexpected at #%d", i)
	}
}
