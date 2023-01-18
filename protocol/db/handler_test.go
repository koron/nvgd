package db

import (
	"testing"
)

func TestBadQuery(t *testing.T) {
	f := func(s string, expected bool) {
		if reBadQuery.MatchString(s) != expected {
			t.Errorf("reBadQuery.MatchString(%q) should be %t", s, expected)
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
	// variations
	f("insert", true)
	f("UpDaTe", true)
	f(" DELETE", true)
	f(" CrEaTe", true)
	f(" DROP", true)
	f(" ALTer", true)
	// inhibits
	f("SELECT * FROM USERS", false)
	f("SELECT * FROM UPDATES", false)
}

func TestHasLimit(t *testing.T) {
	for i, tc := range []struct {
		query string
		want  bool
	}{
		{``, false},
		{`SELECT gid, COUNT(*) FROM user`, true}, // SELECT with COUNT
		{`select gid, COUNT(*) from user`, true}, // case insensitive
		{`SELECT * FROM user LIMIT 10000`, true}, // LIMIT with number
		{`LIMIT 10000`, true},                    // only LIMIT with number
		{`SELECT * FROM user LIMIT abc`, false},  // LIMIT without number
	} {
		got := hasLimit(tc.query)
		if got != tc.want {
			t.Errorf("unexpected hasLimit return: want=%t got=%t:#%d: %s", tc.want, got, i, tc.query)
		}
	}
}
