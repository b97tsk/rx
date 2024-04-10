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
	Source   Observable[T]
	Duration time.Duration
}

func (obs delayObservable[T]) Subscribe(c Context, o Observer[T]) {
	c, cancel := c.WithCancel()
	o = o.DoOnTermination(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Queue    struct {
			sync.Mutex
			queue.Queue[Pair[time.Time, T]]
		}
		Worker struct {
			sync.WaitGroup
		}
	}

	x.Context.Store(c.Context)

	var startWorker func(time.Duration)

	startWorker = func(timeout time.Duration) {
		w, cancelw := c.WithCancel()

		x.Context.Store(w.Context)
		x.Worker.Add(1)

		Timer(timeout).Subscribe(w, func(n Notification[time.Time]) {
			switch n.Kind {
			case KindNext:
				x.Queue.Lock()

				done := w.Done()

				for {
					select {
					default:
					case <-done:
						old := x.Context.Swap(sentinel)

						x.Queue.Init()
						x.Queue.Unlock()

						if old != sentinel {
							o.Error(w.Err())
						}

						return
					}

					n := x.Queue.Front()

					if d := time.Until(n.Key); d > 0 {
						x.Queue.Unlock()
						startWorker(d)
						return
					}

					x.Queue.Pop()

					x.Queue.Unlock()
					o.Next(n.Value)
					x.Queue.Lock()

					if x.Queue.Len() == 0 {
						swapped := x.Context.CompareAndSwap(w.Context, c.Context)

						x.Queue.Unlock()

						if swapped && x.Complete.Load() && x.Context.CompareAndSwap(c.Context, sentinel) {
							o.Complete()
						}

						return
					}
				}

			case KindError:
				defer x.Worker.Done()

				cancelw()

				x.Queue.Lock()
				old := x.Context.Swap(sentinel)
				x.Queue.Init()
				x.Queue.Unlock()

				if old != sentinel {
					o.Error(n.Error)
				}

			case KindComplete:
				cancelw()
				x.Worker.Done()
			}
		})
	}

	obs.Source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			x.Queue.Lock()

			ctx := x.Context.Load()
			if ctx != sentinel {
				x.Queue.Push(NewPair(time.Now().Add(obs.Duration), n.Value))
			}

			x.Queue.Unlock()

			if ctx == c.Context {
				startWorker(obs.Duration)
			}

		case KindError:
			old := x.Context.Swap(sentinel)

			cancel()
			x.Worker.Wait()
			x.Queue.Init()

			if old != sentinel {
				o.Emit(n)
			}

		case KindComplete:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(c.Context, sentinel) {
				o.Emit(n)
			}
		}
	})
}
