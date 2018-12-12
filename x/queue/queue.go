package queue

const minCapacity = 4

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
	q.growIfFull()
	q.buf[q.tail] = x
	q.tail = (q.tail + 1) % len(q.buf)
	q.length++
}

// PopFront removes and returns the first element. It panics if the queue
// is empty.
func (q *Queue) PopFront() interface{} {
	if q.length == 0 {
		panic("queue: PopFront() called on empty queue")
	}
	x := q.buf[q.head]
	q.buf[q.head] = nil
	q.head = (q.head + 1) % len(q.buf)
	q.length--
	q.shrinkIfSparse()
	return x
}

// At returns the i-th element in the queue. It panics if the queue is empty
// or i is not between 0 and Len()-1.
func (q *Queue) At(i int) interface{} {
	if i < 0 || i >= q.length {
		panic("queue: At() called with index out of range")
	}
	return q.buf[(q.head+i)%len(q.buf)]
}

// Front returns the first element. It panics if the queue is empty.
func (q *Queue) Front() interface{} {
	if q.length == 0 {
		panic("queue: Front() called on empty queue")
	}
	return q.buf[q.head]
}

// Back returns the last element. It panics if the queue is empty.
func (q *Queue) Back() interface{} {
	if q.length == 0 {
		panic("queue: Back() called on empty queue")
	}
	return q.buf[(q.tail+len(q.buf)-1)%len(q.buf)]
}

func (q *Queue) growIfFull() {
	if q.length == len(q.buf) {
		q.resize()
	}
}

func (q *Queue) shrinkIfSparse() {
	if q.length == len(q.buf)/4 && len(q.buf) > minCapacity {
		q.resize()
	}
}

func (q *Queue) resize() {
	size := q.length * 2
	if size < minCapacity {
		size = minCapacity
	}
	buf := make([]interface{}, size)
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
