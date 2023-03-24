package rx

import (
	"context"
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

func (obs delayObservable[T]) Subscribe(ctx context.Context, sink Observer[T]) {
	source, cancel := context.WithCancel(ctx)

	sink = sink.OnLastNotification(cancel)

	var x struct {
		Context  atomic.Value
		Complete atomic.Bool
		Queue    struct {
			sync.Mutex
			queue.Queue[Pair[time.Time, T]]
		}
		Worker struct {
			sync.WaitGroup
			Cancel context.CancelFunc
		}
	}

	x.Context.Store(source)

	var startWorker func(time.Duration)

	startWorker = func(timeout time.Duration) {
		worker, cancel := context.WithCancel(source)
		x.Context.Store(worker)

		x.Worker.Cancel = cancel

		x.Worker.Add(1)

		Timer(timeout).Subscribe(worker, func(n Notification[time.Time]) {
			switch {
			case n.HasValue:
				x.Queue.Lock()

				done := worker.Done()

				for {
					select {
					default:
					case <-done:
						swapped := x.Context.CompareAndSwap(worker, sentinel)

						x.Queue.Init()
						x.Queue.Unlock()

						if swapped {
							sink.Error(worker.Err())
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
					sink.Next(n.Value)
					x.Queue.Lock()

					if x.Queue.Len() == 0 {
						swapped := x.Context.CompareAndSwap(worker, source)

						x.Queue.Unlock()

						if swapped && x.Complete.Load() && x.Context.CompareAndSwap(source, sentinel) {
							sink.Complete()
						}

						return
					}
				}

			case n.HasError:
				x.Queue.Lock()

				swapped := x.Context.CompareAndSwap(worker, sentinel)

				x.Queue.Init()
				x.Queue.Unlock()

				if swapped {
					sink.Error(n.Error)
				}

			default:
				break
			}

			cancel()
			x.Worker.Done()
		})
	}

	obs.Source.Subscribe(source, func(n Notification[T]) {
		switch {
		case n.HasValue:
			x.Queue.Lock()

			ctx := x.Context.Load()
			if ctx != sentinel {
				x.Queue.Push(NewPair(time.Now().Add(obs.Duration), n.Value))
			}

			x.Queue.Unlock()

			if ctx == source {
				startWorker(obs.Duration)
			}

		case n.HasError:
			ctx := x.Context.Swap(source)

			if x.Worker.Cancel != nil {
				x.Worker.Cancel()
				x.Worker.Wait()
			}

			if x.Context.Swap(sentinel) != sentinel && ctx != sentinel {
				sink(n)
			}

			x.Queue.Init()

		default:
			x.Complete.Store(true)

			if x.Context.CompareAndSwap(source, sentinel) {
				sink(n)
			}
		}
	})
}
