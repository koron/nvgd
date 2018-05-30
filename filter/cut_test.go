package filter

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func checkCut(t *testing.T, cutSrc, cutfmt, expected string) {
	t.Helper()
	sels, err := toCutSelector(cutfmt)
	if err != nil {
		t.Errorf("invalid cut format: %s", err)
		return
	}

	rb := ioutil.NopCloser(bytes.NewReader([]byte(cutSrc)))
	cut := NewCut(rb, []byte("\t"), sels)
	defer func() {
		err := cut.Close()
		if err != nil {
			t.Errorf("failed to close: %s", err)
		}
	}()

	b, err := ioutil.ReadAll(cut)
	if err != nil {
		t.Errorf("failed to read: %s", err)
		return
	}
	act := string(b)
	if act != expected {
		t.Errorf("cut returns unexpected\nactual:%q\nexpected:%q\n", act, expected)
	}
}

func TestCutSelector(t *testing.T) {
	const src = "A\tB\tC\tD\tE\tF\tG\tH\tI\tJ\tK\tL\tM\tN\tO\tP\tQ\tR\tS\tT\tU\tV\tW\tX\tY\tZ\na\tb\tc\td\te\tf\tg\th\ti\tj\tk\tl\tm\tn\to\tp\tq\tr\ts\tt\tu\tv\tw\tx\ty\tz"

	checkCut(t, src, "1", "A\na")
	checkCut(t, src, "12", "L\nl")

	checkCut(t, src, "11-15", "K\tL\tM\tN\tO\nk\tl\tm\tn\to")
	checkCut(t, src, "21-25", "U\tV\tW\tX\tY\nu\tv\tw\tx\ty")
	// reverse range
	checkCut(t, src, "13-11", "M\tL\tK\nm\tl\tk")

	checkCut(t, src, "24-", "X\tY\tZ\nx\ty\tz")
	checkCut(t, src, "-3", "A\tB\tC\na\tb\tc")

	// combinations

	checkCut(t, src, "1,5", "A\tE\na\te")
	checkCut(t, src, "12,26", "L\tZ\nl\tz")

	// TODO: need complex combinations
}

func TestCutEmpty(t *testing.T) {
	// empty lines are at top, middle and bottom.
	const src = `
A1	B1	C1
A2	B2	C2

A4	B4	C4
A5	B5	C5

`
	checkCut(t, src, "1", "\nA1\nA2\n\nA4\nA5\n\n")
	checkCut(t, src, "2", "\nB1\nB2\n\nB4\nB5\n\n")
	checkCut(t, src, "3", "\nC1\nC2\n\nC4\nC5\n\n")
}
