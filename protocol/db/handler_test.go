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

func TestExpandName(t *testing.T) {
	h := &Handler{}
	f := func(driver, name, dbname string, expected string) {
		actual, err := h.expandName(driver, name, dbname)
		if err != nil {
			t.Fatalf("error: %s\ndriver=%q\nname=%q\ndbname=%q",
				err, driver, name, dbname)
		}
		if actual != expected {
			t.Errorf("unmatch: expected=%q actual=%q\ndriver=%q\nname=%q\ndbname=%q",
				expected, actual, driver, name, dbname)
		}
	}
	f("", "foo", "bar", "foo")
	f("", "foo{{.dbname}}", "bar", "foobar")
	f("", "foo{{.dbname}}baz", "bar", "foobarbaz")
	//f("", "foo{{.name}}baz", "bar", "foobarbaz")
}
