package ringbuf

type Buffer struct {
	d []interface{}
	r int
	n int
}

func New(capacity int) *Buffer {
	if capacity <= 0 {
		panic("capacity must be large than 0")
	}
	return &Buffer{
		d: make([]interface{}, capacity),
	}
}

func (b *Buffer) Put(v interface{}) {
	b.d[(b.r+b.n)%len(b.d)] = v
	if b.n < len(b.d) {
		b.n += 1
	} else {
		b.incR()
	}
}

func (b *Buffer) Get() (interface{}, bool) {
	if b.n <= 0 {
		return nil, false
	}
	v := b.d[b.r]
	b.d[b.r] = nil
	b.incR()
	b.n -= 1
	return v, true
}

func (b *Buffer) Clear() {
	for b.n > 0 {
		b.d[b.r] = nil
		b.incR()
		b.n -= 1
	}
}

func (b *Buffer) incR() int {
	b.r += 1
	if b.r == len(b.d) {
		b.r = 0
	}
	return b.r
}
