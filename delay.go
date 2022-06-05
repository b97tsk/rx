package rx

import (
	"context"
	"time"

	"github.com/b97tsk/rx/internal/critical"
	"github.com/b97tsk/rx/internal/queue"
)

// Delay delays the emission of items from the source Observable by a given
// timeout.
func Delay[T any](d time.Duration) Operator[T, T] {
	return AsOperator(
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
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Queue     queue.Queue[Pair[time.Time, T]]
		Scheduled bool
		Completed bool
	}

	var onTimer Observer[time.Time]

	doSchedule := func(timeout time.Duration) {
		Timer(timeout).Subscribe(ctx, onTimer)
	}

	onTimer = func(n Notification[time.Time]) {
		if !n.HasValue && critical.Enter(&x.Section) {
			if n.HasError {
				critical.Close(&x.Section)

				sink.Error(n.Error)

				return
			}

			x.Scheduled = false

			for x.Queue.Len() > 0 {
				if err := ctx.Err(); err != nil {
					critical.Close(&x.Section)

					sink.Error(err)

					return
				}

				n := x.Queue.Front()

				if d := time.Until(n.Key); d > 0 {
					x.Scheduled = true

					doSchedule(d)

					break
				}

				x.Queue.Pop()

				sink.Next(n.Value)
			}

			if x.Completed && x.Queue.Len() == 0 {
				sink.Complete()
			}

			critical.Leave(&x.Section)
		}
	}

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		if critical.Enter(&x.Section) {
			switch {
			case n.HasValue:
				x.Queue.Push(MakePair(time.Now().Add(obs.Duration), n.Value))

				if !x.Scheduled {
					x.Scheduled = true
					doSchedule(obs.Duration)
				}

				critical.Leave(&x.Section)

			case n.HasError:
				x.Queue.Init()

				critical.Close(&x.Section)

				sink(n)

			default:
				x.Completed = true

				if x.Queue.Len() > 0 {
					critical.Leave(&x.Section)
					break
				}

				critical.Close(&x.Section)

				sink(n)
			}
		}
	})
}
