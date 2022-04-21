// Package ringbuf provides ring buffer.
package ringbuf

// Buffer provides ring buffer.
type Buffer struct {
	d []interface{}
	r int
	n int
}

// New creates a ring buffer with specified capacity.
func New(capacity int) *Buffer {
	if capacity <= 0 {
		panic("capacity must be large than 0")
	}
	return &Buffer{
		d: make([]interface{}, capacity),
	}
}

// Put puts a value.
func (b *Buffer) Put(v interface{}) {
	b.d[(b.r+b.n)%len(b.d)] = v
	if b.n < len(b.d) {
		b.n++
	} else {
		b.incR()
	}
}

// Get retrieves a value.
func (b *Buffer) Get() (interface{}, bool) {
	if b.n <= 0 {
		return nil, false
	}
	v := b.d[b.r]
	b.d[b.r] = nil
	b.incR()
	b.n--
	return v, true
}

// Clear remove all values.
func (b *Buffer) Clear() {
	for b.n > 0 {
		b.d[b.r] = nil
		b.incR()
		b.n--
	}
}

// Empty checks the buffer is empty or not.
func (b *Buffer) Empty() bool {
	return b.n == 0
}

func (b *Buffer) incR() int {
	b.r++
	if b.r == len(b.d) {
		b.r = 0
	}
	return b.r
}
