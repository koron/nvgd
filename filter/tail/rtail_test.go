package tail

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

type nopCloser struct {
	io.ReadSeeker
}

func (nopCloser) Close() error { return nil }

func newRSC(src string) readSeekCloser {
	b := bytes.NewReader([]byte(src))
	return nopCloser{b}
}

func checkRTail(t *testing.T, src string, limit, bufsize int, exp string) {
	rt := NewRTail(newRSC(src), limit, bufsize)
	defer func() {
		err := rt.Close()
		if err != nil {
			t.Errorf("failed to close: %s", err)
		}
	}()
	b, err := ioutil.ReadAll(rt)
	if err != nil {
		t.Errorf("failed to read: %s", err)
	}
	act := string(b)
	if act != exp {
		t.Errorf("rtail returns unexpected\nactual:%q\nexpect:%q\nlimit:%d bufsize:%d", act, exp, limit, bufsize)
	}
}

const rtailSrc = `aaa
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

func TestRTail(t *testing.T) {
	checkRTail(t, `abc`, 10, 4096, `abc`)
	checkRTail(t, rtailSrc, 1, 4096, "jjj\n")
	checkRTail(t, rtailSrc, 1, 5, "jjj\n")
	checkRTail(t, rtailSrc, 3, 5, "hhh\niii\njjj\n")
	checkRTail(t, rtailSrc, 11, 5, rtailSrc)
	checkRTail(t, rtailSrc[:len(rtailSrc)-1], 1, 5, "jjj")
	checkRTail(t, rtailSrc[:len(rtailSrc)-1], 3, 5, "hhh\niii\njjj")
}
