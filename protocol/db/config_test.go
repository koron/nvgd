package db

import (
	"net/url"
	"testing"
)

func TestExpandName(t *testing.T) {
	f := func(driver, name, dbname string, expected string) {
		p := &Param{
			Driver: driver,
			Name:   name,
		}
		actual, err := p.expandName(dbname)
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

func TestParam(t *testing.T) {
	p := Param{
		Driver: "mysql",
		Name:   "nvgd:nvgd@/{{.dbname}}",

		MultipleDatabase: true,
	}
	ok := func(name string, exp string) {
		t.Helper()
		act, err := p.expandName(name)
		if err != nil {
			t.Fatalf("failed to expandName: %s", err)
		}
		if act != exp {
			t.Errorf("expandName returns wrong: %q (expected %q)", act, exp)
		}
	}
	ok("test", `nvgd:nvgd@/test`)
	ok("foo", `nvgd:nvgd@/foo`)
	ok("", `nvgd:nvgd@/{{.dbname}}`)
}

func TestExtractNames(t *testing.T) {
	ok := func(s, expName, expDbn string) {
		u, err := url.Parse(s)
		if err != nil {
			t.Fatalf("failed to parse url: %s", err)
		}
		name, dbn := extractNames(u)
		if name != expName {
			t.Errorf("unexpected name: %q (expected %q)", name, expName)
		}
		if dbn != expDbn {
			t.Errorf("unexpected dbname: %q (expected %q)", dbn, expDbn)
		}
	}
	ok(`db://db_single/select 1`, `db_single`, ``)
	ok(`db://foo@db_multi/select 1`, `db_multi`, `foo`)
	ok(`db://test@db_multi/select 1`, `db_multi`, `test`)
	ok(`db://@db_multi/select 1`, `db_multi`, ``)
}
