package rx

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/b97tsk/rx/internal/queue"
)

// Delay postpones each emission of values from the source Observable
// by a given duration.
func Delay[T any](d time.Duration) Operator[T, T] {
	return NewOperator(
		func(source Observable[T]) Observable[T] {
			return delayObservable[T]{source, d}.Subscribe
		},
	)
}

type delayObservable[T any] struct {
	source   Observable[T]
	duration time.Duration
}

func (ob delayObservable[T]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		context  atomic.Value
		complete atomic.Bool
		buffer   struct {
			sync.Mutex
			queue.Queue[Pair[time.Time, T]]
		}
		worker struct {
			sync.WaitGroup
		}
	}

	x.context.Store(c.Context)

	var startWorker func(time.Duration)

	startWorker = func(timeout time.Duration) {
		w, cancelw := c.WithCancel()

		x.context.Store(w.Context)
		x.worker.Add(1)

		Timer(timeout).Subscribe(w, func(n Notification[time.Time]) {
			switch n.Kind {
			case KindNext:
				x.buffer.Lock()

				done := w.Done()

				for {
					select {
					default:
					case <-done:
						old := x.context.Swap(sentinel)

						x.buffer.Init()
						x.buffer.Unlock()

						if old != sentinel {
							o.Error(w.Cause())
						}

						return
					}

					n := x.buffer.Front()

					if d := time.Until(n.Key); d > 0 {
						x.buffer.Unlock()
						startWorker(d)
						return
					}

					x.buffer.Pop()

					x.buffer.Unlock()
					o.Next(n.Value)
					x.buffer.Lock()

					if x.buffer.Len() == 0 {
						swapped := x.context.CompareAndSwap(w.Context, c.Context)

						x.buffer.Unlock()

						if swapped && x.complete.Load() && x.context.CompareAndSwap(c.Context, sentinel) {
							o.Complete()
						}

						return
					}
				}

			case KindError:
				defer x.worker.Done()

				cancelw()

				x.buffer.Lock()
				old := x.context.Swap(sentinel)
				x.buffer.Init()
				x.buffer.Unlock()

				if old != sentinel {
					o.Error(n.Error)
				}

			case KindComplete:
				cancelw()
				x.worker.Done()
			}
		})
	}

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.buffer.Lock()

			ctx := x.context.Load()
			if ctx != sentinel {
				x.buffer.Push(NewPair(time.Now().Add(ob.duration), n.Value))
			}

			x.buffer.Unlock()

			if ctx == c.Context {
				startWorker(ob.duration)
			}

		case KindError:
			old := x.context.Swap(sentinel)

			cancel()
			x.worker.Wait()
			x.buffer.Init()

			if old != sentinel {
				o.Emit(n)
			}

		case KindComplete:
			x.complete.Store(true)

			if x.context.CompareAndSwap(c.Context, sentinel) {
				o.Emit(n)
			}
		}
	})
}
