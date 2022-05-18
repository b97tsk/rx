package queue

const smallSize = 16

// Queue represents a single instance of the queue data structure. The zero
// value for a Queue is an empty queue ready to use.
type Queue[E any] struct {
	head, tail []E
}

// Init initializes or clears q.
func (q *Queue[E]) Init() {
	q.head, q.tail = nil, nil
}

// Cap returns the capacity of the internal buffer. If Cap() equals to Len(),
// new Push(x) causes the internal buffer to grow.
func (q *Queue[E]) Cap() int {
	return cap(q.tail)
}

// Len returns the number of elements currently stored in q.
func (q *Queue[E]) Len() int {
	return len(q.head) + len(q.tail)
}

// Push inserts an element at the end of q.
func (q *Queue[E]) Push(x E) {
	if q.Len() == q.Cap() { // Grow if full.
		buf := append(append(q.head, q.tail...), x)
		q.head, q.tail = buf, buf[:0]

		return
	}

	if len(q.head) < cap(q.head) {
		q.head = append(q.head, x)
	} else {
		q.tail = append(q.tail, x)
	}
}

// Pop removes and returns the first element. It panics if q is empty.
func (q *Queue[E]) Pop() E {
	if len(q.head) > 0 {
		x := q.head[0]

		var zero E
		q.head[0] = zero
		q.head = q.head[1:]

		if cap(q.head) == 0 {
			q.head = q.tail
			q.tail = q.tail[:0]
		}

		if n, m := q.Len(), q.Cap(); n == m>>2 && m > smallSize { // Shrink if sparse.
			buf := make([]E, n<<1)
			q.head, q.tail = buf[:q.CopyTo(buf)], buf[:0]
		}

		return x
	}

	panic("queue: Pop called on empty Queue")
}

// At returns the i-th element in q. It panics if i is out of range.
func (q *Queue[E]) At(i int) E {
	if i >= 0 {
		headsize := len(q.head)

		if i < headsize {
			return q.head[i]
		}

		i -= headsize

		if i < len(q.tail) {
			return q.tail[i]
		}
	}

	panic("queue: At called with index out of range")
}

// Front returns the first element. It panics if q is empty.
func (q *Queue[E]) Front() E {
	if len(q.head) > 0 {
		return q.head[0]
	}

	panic("queue: Front called on empty Queue")
}

// Back returns the last element. It panics if q is empty.
func (q *Queue[E]) Back() E {
	if n := len(q.tail); n > 0 {
		return q.tail[n-1]
	}

	if n := len(q.head); n > 0 {
		return q.head[n-1]
	}

	panic("queue: Back called on empty Queue")
}

// CopyTo copies elements of q into a destination slice and returns
// the number of elements copied, which will be the minimum of q.Len()
// and len(dst).
func (q *Queue[E]) CopyTo(dst []E) int {
	n := copy(dst, q.head)
	return n + copy(dst[n:], q.tail)
}

// Clone clones q.
func (q *Queue[E]) Clone() Queue[E] {
	var buf []E

	if q.head != nil {
		buf = make([]E, q.Len(), q.Cap())
		q.CopyTo(buf)
	}

	return Queue[E]{buf, buf[:0]}
}
