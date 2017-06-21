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
