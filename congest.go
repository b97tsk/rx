package rx

import "github.com/b97tsk/rx/internal/queue"

// CongestBlock mirrors the source Observable, buffering emissions
// if the source emits too fast, and blocking the source if the buffer is
// full.
func CongestBlock[T any](capacity int) Operator[T, T] {
	return Channelize(
		func(upstream <-chan Notification[T], downstream chan<- Notification[T]) {
			if capacity < 1 {
				panic("CongestBlock: capacity < 1")
			}

			var buf queue.Queue[Notification[T]]

			var complete bool

			for {
				var (
					in   <-chan Notification[T]
					out  chan<- Notification[T]
					outv Notification[T]
				)

				length := buf.Len()

				if length < capacity {
					in = upstream
				}

				if length > 0 {
					out, outv = downstream, buf.Front()
				}

				select {
				case n := <-in:
					switch n.Kind {
					case KindNext:
						buf.Push(n)
					case KindError:
						downstream <- n
						return
					case KindComplete:
						complete = true

						if buf.Len() == 0 {
							downstream <- n
							return
						}
					}
				case out <- outv:
					buf.Pop()

					if complete && buf.Len() == 0 {
						downstream <- Complete[T]()
						return
					}
				}
			}
		},
	)
}

// CongestDropLatest mirrors the source Observable, buffering emissions
// if the source emits too fast, and dropping emissions if the buffer is
// full.
func CongestDropLatest[T any](capacity int) Operator[T, T] {
	return Channelize(
		func(upstream <-chan Notification[T], downstream chan<- Notification[T]) {
			if capacity < 1 {
				panic("CongestDropLatest: capacity < 1")
			}

			var buf queue.Queue[Notification[T]]

			var complete bool

			for {
				var (
					in   = upstream
					out  chan<- Notification[T]
					outv Notification[T]
				)

				if buf.Len() != 0 {
					out, outv = downstream, buf.Front()
				}

				select {
				case n := <-in:
					switch n.Kind {
					case KindNext:
						if buf.Len() < capacity {
							buf.Push(n)
						}
					case KindError:
						downstream <- n
						return
					case KindComplete:
						complete = true

						if buf.Len() == 0 {
							downstream <- n
							return
						}
					}
				case out <- outv:
					buf.Pop()

					if complete && buf.Len() == 0 {
						downstream <- Complete[T]()
						return
					}
				}
			}
		},
	)
}

// CongestDropOldest mirrors the source Observable, buffering emissions
// if the source emits too fast, and dropping oldest emissions from
// the buffer if it is full.
func CongestDropOldest[T any](capacity int) Operator[T, T] {
	return Channelize(
		func(upstream <-chan Notification[T], downstream chan<- Notification[T]) {
			if capacity < 1 {
				panic("CongestDropOldest: capacity < 1")
			}

			var buf queue.Queue[Notification[T]]

			var complete bool

			for {
				var (
					in   = upstream
					out  chan<- Notification[T]
					outv Notification[T]
				)

				if buf.Len() != 0 {
					out, outv = downstream, buf.Front()
				}

				select {
				case n := <-in:
					switch n.Kind {
					case KindNext:
						if buf.Len() == capacity {
							buf.Pop()
						}

						buf.Push(n)
					case KindError:
						downstream <- n
						return
					case KindComplete:
						complete = true

						if buf.Len() == 0 {
							downstream <- n
							return
						}
					}
				case out <- outv:
					buf.Pop()

					if complete && buf.Len() == 0 {
						downstream <- Complete[T]()
						return
					}
				}
			}
		},
	)
}

// CongestError mirrors the source Observable, buffering emissions
// if the source emits too fast, and terminating the subscription with
// a notification of ErrBufferOverflow if the buffer is full.
func CongestError[T any](capacity int) Operator[T, T] {
	return Channelize(
		func(upstream <-chan Notification[T], downstream chan<- Notification[T]) {
			if capacity < 1 {
				panic("CongestError: capacity < 1")
			}

			var buf queue.Queue[Notification[T]]

			var complete bool

			for {
				var (
					in   = upstream
					out  chan<- Notification[T]
					outv Notification[T]
				)

				if buf.Len() != 0 {
					out, outv = downstream, buf.Front()
				}

				select {
				case n := <-in:
					switch n.Kind {
					case KindNext:
						if buf.Len() == capacity {
							downstream <- Error[T](ErrBufferOverflow)
							return
						}

						buf.Push(n)
					case KindError:
						downstream <- n
						return
					case KindComplete:
						complete = true

						if buf.Len() == 0 {
							downstream <- n
							return
						}
					}
				case out <- outv:
					buf.Pop()

					if complete && buf.Len() == 0 {
						downstream <- Complete[T]()
						return
					}
				}
			}
		},
	)
}
