// Package ringbuf provides ring buffer.
package ringbuf

// Buffer provides ring buffer.
type Buffer[T any] struct {
	d []T
	r int
	n int
}

// New creates a ring buffer with specified capacity.
func New[T any](capacity int) *Buffer[T] {
	if capacity <= 0 {
		panic("capacity must be large than 0")
	}
	return &Buffer[T]{
		d: make([]T, capacity),
	}
}

// Put puts a value.
func (b *Buffer[T]) Put(v T) {
	b.d[(b.r+b.n)%len(b.d)] = v
	if b.n < len(b.d) {
		b.n++
	} else {
		b.incR()
	}
}

// Get retrieves a value.
func (b *Buffer[T]) Get() (T, bool) {
	var zero T
	if b.n <= 0 {
		return zero, false
	}
	v := b.d[b.r]
	b.d[b.r] = zero
	b.incR()
	b.n--
	return v, true
}

// Clear remove all values.
func (b *Buffer[T]) Clear() {
	var zero T
	for b.n > 0 {
		b.d[b.r] = zero
		b.incR()
		b.n--
	}
}

// Empty checks the buffer is empty or not.
func (b *Buffer[T]) Empty() bool {
	return b.n == 0
}

func (b *Buffer[T]) incR() int {
	b.r++
	if b.r == len(b.d) {
		b.r = 0
	}
	return b.r
}

// Len returns number of valid items in the ringbuf
func (b *Buffer[T]) Len() int {
	return b.n
}

// Peek peeks checks a n'th value in ringbuf without removing the value.
func (b *Buffer[T]) Peek(n int) T {
	var zero T
	if n < 0 || n >= b.n {
		return zero
	}
	return b.d[(b.r+n)%len(b.d)]
}
