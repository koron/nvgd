package db

import "testing"

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
