package queue

const smallSize = 16

// Queue represents a single instance of the queue data structure. The zero
// value for a Queue is an empty queue ready to use.
type Queue struct {
	buf  []interface{}
	head []interface{}
}

// Init initializes or clears the Queue.
func (q *Queue) Init() {
	q.buf, q.head = nil, nil
}

// Cap returns the capacity of the internal buffer. If Cap() equals to Len(),
// new Push(x) causes the internal buffer to grow.
func (q *Queue) Cap() int {
	return cap(q.buf)
}

// Len returns the number of elements currently stored in the Queue.
func (q *Queue) Len() int {
	return len(q.head) + len(q.buf)
}

// Push inserts an element at the end of the Queue.
func (q *Queue) Push(x interface{}) {
	if ql, qc := q.Len(), q.Cap(); ql == qc { // Grow if full.
		buf := append(q.buf[:qc], nil)
		q.setbuf(buf[:cap(buf)])
	}

	if len(q.head) < cap(q.head) {
		q.head = append(q.head, x)
	} else {
		q.buf = append(q.buf, x)
	}
}

// Pop removes and returns the first element. It panics if the Queue is empty.
func (q *Queue) Pop() interface{} {
	if len(q.head) > 0 {
		x := q.head[0]

		q.head[0] = nil
		q.head = q.head[1:]

		if cap(q.head) == 0 {
			q.head = q.buf
			q.buf = q.buf[:0]
		}

		if ql, qc := q.Len(), q.Cap(); ql == qc>>2 && qc > smallSize { // Shrink if sparse.
			q.setbuf(make([]interface{}, ql<<1))
		}

		return x
	}

	panic("queue: Pop called on empty Queue")
}

// At returns the i-th element in the Queue. It panics if i is out of range.
func (q *Queue) At(i int) interface{} {
	if i >= 0 {
		headsize := len(q.head)

		if i < headsize {
			return q.head[i]
		}

		i -= headsize

		if i < len(q.buf) {
			return q.buf[i]
		}
	}

	panic("queue: At called with index out of range")
}

// Front returns the first element. It panics if the Queue is empty.
func (q *Queue) Front() interface{} {
	if len(q.head) > 0 {
		return q.head[0]
	}

	panic("queue: Front called on empty Queue")
}

// Back returns the last element. It panics if the Queue is empty.
func (q *Queue) Back() interface{} {
	if n := len(q.buf); n > 0 {
		return q.buf[n-1]
	}

	if n := len(q.head); n > 0 {
		return q.head[n-1]
	}

	panic("queue: Back called on empty Queue")
}

// CopyTo copies elements of the Queue into a destination slice. CopyTo
// returns the number of elements copied, which will be the minimum of
// q.Len() and len(dst).
func (q *Queue) CopyTo(dst []interface{}) int {
	n := copy(dst, q.head)
	return n + copy(dst[n:], q.buf)
}

// Clone clones the Queue.
func (q *Queue) Clone() Queue {
	var buf []interface{}

	if q.buf != nil {
		buf = make([]interface{}, q.Len(), q.Cap())
		q.CopyTo(buf)
	}

	return Queue{buf[:0], buf}
}

func (q *Queue) setbuf(buf []interface{}) {
	q.buf, q.head = buf[:0], buf[:q.CopyTo(buf)]
}
