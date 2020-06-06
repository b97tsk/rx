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

// Init initializes or clears the queue.
func (q *Queue) Init() {
	*q = Queue{}
}

// Cap returns the capacity of the internal buffer. If Cap() equals to Len(),
// new PushBack(x) causes the internal buffer to grow.
func (q *Queue) Cap() int {
	return len(q.buf)
}

// Len returns the number of elements currently stored in the queue.
func (q *Queue) Len() int {
	return q.length
}

// PushBack inserts an element at the end of the queue.
func (q *Queue) PushBack(x interface{}) {
	if q.length == len(q.buf) { // Grow if full.
		buf := append(q.buf, x)
		q.setbuf(buf[:cap(buf)])
	}
	q.buf[q.tail] = x
	q.tail = (q.tail + 1) % len(q.buf)
	q.length++
}

// PopFront removes and returns the first element. It panics if the queue
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

// At returns the i-th element in the queue. It panics if the queue is empty
// or i is not between 0 and Len()-1.
func (q *Queue) At(i int) interface{} {
	return q.buf[(q.head+i)%len(q.buf)]
}

// Front returns the first element. It panics if the queue is empty.
func (q *Queue) Front() interface{} {
	return q.buf[q.head]
}

// Back returns the last element. It panics if the queue is empty.
func (q *Queue) Back() interface{} {
	max := len(q.buf)
	return q.buf[(q.tail+max-1)%max]
}

func (q *Queue) setbuf(buf []interface{}) {
	if q.head < q.tail {
		copy(buf, q.buf[q.head:q.tail])
	} else {
		n := copy(buf, q.buf[q.head:])
		copy(buf[n:], q.buf[:q.tail])
	}
	q.buf = buf
	q.head = 0
	q.tail = q.length
}
