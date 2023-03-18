package tail

import (
	"bytes"
	"io"
	"testing"
)

func checkTail(t *testing.T, src string, limit int, exp string) {
	rb := io.NopCloser(bytes.NewReader([]byte(src)))
	rt := NewTail(rb, limit)
	defer func() {
		err := rt.Close()
		if err != nil {
			t.Errorf("failed to close: %s", err)
		}
	}()
	b, err := io.ReadAll(rt)
	if err != nil {
		t.Errorf("failed to read: %s", err)
	}
	act := string(b)
	if act != exp {
		t.Errorf("tail returns unexpected\nactual:%q\nexpect:%q\nlimit:%d", act, exp, limit)
	}
}

const tailSrc = `aaa
bbb
ccc
ddd
eee
fff
ggg
hhh
iii
jjj
`

func TestTail(t *testing.T) {
	checkTail(t, `abc`, 10, `abc`)
	checkTail(t, tailSrc, 1, "jjj\n")
	checkTail(t, tailSrc, 3, "hhh\niii\njjj\n")
	checkTail(t, tailSrc, 11, tailSrc)
	checkTail(t, tailSrc[:len(tailSrc)-1], 1, "jjj")
	checkTail(t, tailSrc[:len(tailSrc)-1], 3, "hhh\niii\njjj")
}
