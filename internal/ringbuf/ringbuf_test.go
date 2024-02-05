package ringbuf

import "testing"

func TestRing(t *testing.T) {
	check := func(n int, in []int, exp []int) {
		b := New[int](n)
		for _, v := range in {
			b.Put(v)
		}
		for i, v := range exp {
			w, ok := b.Get()
			if !ok {
				t.Errorf("Get() failed: i=%d n=%d in=%v", i, n, in)
				return
			}
			if w != v {
				t.Errorf("Get() mismatch: i=%d w=%v(!=%d) n=%d in=%v", i, w, v, n, in)
				return
			}
		}
		w, ok := b.Get()
		if ok {
			t.Errorf("more value: w=%v n=%d in=%v", w, n, in)
		}
	}

	for _, d := range []struct {
		n   int
		in  []int
		exp []int
	}{
		// length: 1
		{1, []int{}, []int{}},
		{1, []int{1}, []int{1}},
		{1, []int{1, 2}, []int{2}},
		{1, []int{1, 2, 3}, []int{3}},
		{1, []int{1, 2, 3, 4}, []int{4}},

		// length: 2
		{2, []int{}, []int{}},
		{2, []int{1}, []int{1}},
		{2, []int{1, 2}, []int{1, 2}},
		{2, []int{1, 2, 3}, []int{2, 3}},
		{2, []int{1, 2, 3, 4}, []int{3, 4}},

		// length: 5
		{5, []int{}, []int{}},
		{5, []int{9}, []int{9}},
		{5, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}},
		{5, []int{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{4, 5, 6, 7, 8}},
	} {
		check(d.n, d.in, d.exp)
	}
}

func TestRingExtra(t *testing.T) {
	rb := New[int](5)
	if want, got := 0, rb.Len(); got != want {
		t.Fatalf("invalid length: want=%d got=%d", want, got)
	}
	for i := 1; i <= 4; i++ {
		rb.Put(i)
		if want, got := i, rb.Len(); got != want {
			t.Fatalf("invalid length: want=%d got=%d", want, got)
		}
	}
	for i := 5; i <= 8; i++ {
		rb.Put(i)
		if want, got := 5, rb.Len(); got != want {
			t.Fatalf("invalid length: want=%d got=%d", want, got)
		}
	}

	// test ringbuf.Peek
	if want, got := 4, rb.Peek(0); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", 0, want, got)
	}
	if want, got := 5, rb.Peek(1); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", 1, want, got)
	}
	if want, got := 6, rb.Peek(2); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", 2, want, got)
	}
	if want, got := 7, rb.Peek(3); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", 3, want, got)
	}
	if want, got := 8, rb.Peek(4); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", 4, want, got)
	}

	// ringbuf.Peek with values of out
	if want, got := 0, rb.Peek(6); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", 6, want, got)
	}
	if want, got := 0, rb.Peek(-1); got != want {
		t.Errorf("Peek(%d) returns unexpected: want=%d got=%d", -1, want, got)
	}

	// Clear and Empty
	if want, got := false, rb.Empty(); want != got {
		t.Errorf("Empty() returns unexpected: want=%t got=%t", want, got)
	}
	rb.Clear()
	if want, got := true, rb.Empty(); want != got {
		t.Errorf("Empty() returns unexpected: want=%t got=%t", want, got)
	}
	// Len() should be 0 after Clean()
	if want, got := 0, rb.Len(); want != got {
		t.Errorf("Len() returns unexpected: want=%d got=%d", want, got)
	}
}
