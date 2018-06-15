package ringbuf

import "testing"

func TestRing(t *testing.T) {
	check := func(n int, in []int, exp []int) {
		b := New(n)
		for _, v := range in {
			b.Put(v)
		}
		for i, v := range exp {
			w, ok := b.Get()
			if !ok {
				t.Errorf("Get() failed: i=%d n=%d in=%v", i, n, in)
				return
			}
			if w.(int) != v {
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
