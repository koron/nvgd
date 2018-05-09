package filter

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func checkCut(t *testing.T, cutfmt, expected string) {
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

const cutSrc = "A\tB\tC\tD\tE\tF\tG\tH\tI\tJ\tK\tL\tM\tN\tO\tP\tQ\tR\tS\tT\tU\tV\tW\tX\tY\tZ\na\tb\tc\td\te\tf\tg\th\ti\tj\tk\tl\tm\tn\to\tp\tq\tr\ts\tt\tu\tv\tw\tx\ty\tz"

func TestCutSelector(t *testing.T) {
	checkCut(t, "1", "A\na")
	checkCut(t, "12", "L\nl")

	checkCut(t, "11-15", "K\tL\tM\tN\tO\nk\tl\tm\tn\to")
	checkCut(t, "21-25", "U\tV\tW\tX\tY\nu\tv\tw\tx\ty")
	// reverse range
	checkCut(t, "13-11", "M\tL\tK\nm\tl\tk")

	checkCut(t, "24-", "X\tY\tZ\nx\ty\tz")
	checkCut(t, "-3", "A\tB\tC\na\tb\tc")

	// combinations

	checkCut(t, "1,5", "A\tE\na\te")
	checkCut(t, "12,26", "L\tZ\nl\tz")

	// TODO: need complex combinations
}
