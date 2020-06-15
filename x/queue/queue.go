package queue

const smallSize = 8

// Queue represents a single instance of the queue data structure. The zero
// value for a Queue is an empty queue ready to use.
type Queue struct {
	buf    []interface{}
	head   int
	tail   int
	length int
}

// Init initializes or clears the Queue.
func (q *Queue) Init() {
	*q = Queue{}
}

// Cap returns the capacity of the internal buffer. If Cap() equals to Len(),
// new PushBack(x) causes the internal buffer to grow.
func (q *Queue) Cap() int {
	return len(q.buf)
}

// Len returns the number of elements currently stored in the Queue.
func (q *Queue) Len() int {
	return q.length
}

// PushBack inserts an element at the end of the Queue.
func (q *Queue) PushBack(x interface{}) {
	if q.length == len(q.buf) { // Grow if full.
		buf := append(q.buf, x)
		q.setbuf(buf[:cap(buf)])
	}
	q.buf[q.tail] = x
	q.tail = (q.tail + 1) % len(q.buf)
	q.length++
}

// PopFront removes and returns the first element. It panics if the Queue
// is empty.
func (q *Queue) PopFront() interface{} {
	x := q.buf[q.head]
	q.buf[q.head] = nil
	max := len(q.buf)
	q.head = (q.head + 1) % max
	q.length--
	if q.length == max/4 && max > smallSize { // Shrink if sparse.
		q.setbuf(make([]interface{}, q.length*2))
	}
	return x
}

// At returns the i-th element in the Queue. It panics if i is out of range.
func (q *Queue) At(i int) interface{} {
	return q.buf[(q.head+i)%len(q.buf)]
}

// Front returns the first element. It panics if the Queue is empty.
func (q *Queue) Front() interface{} {
	return q.buf[q.head]
}

// Back returns the last element. It panics if the Queue is empty.
func (q *Queue) Back() interface{} {
	max := len(q.buf)
	return q.buf[(q.tail+max-1)%max]
}

// Copy copies elements of the Queue into a destination slice. Copy returns
// the number of elements copied, which will be the minimum of q.Len() and
// len(dst).
func (q *Queue) Copy(dst []interface{}) int {
	if q.head < q.tail {
		return copy(dst, q.buf[q.head:q.tail])
	}
	n := copy(dst, q.buf[q.head:])
	return n + copy(dst[n:], q.buf[:q.tail])
}

// Clone clones the Queue.
func (q *Queue) Clone() Queue {
	var b []interface{}
	if q.length > 0 {
		b = make([]interface{}, q.length, len(q.buf))
		q.Copy(b)
	}
	return Queue{b, 0, len(b), len(b)}
}

func (q *Queue) setbuf(buf []interface{}) {
	q.Copy(buf)
	q.buf = buf
	q.head = 0
	q.tail = q.length
}
