package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
	"github.com/b97tsk/rx/internal/queue"
)

// Delay delays the emission of items from the source Observable by a given
// timeout.
func Delay(d time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return delayObservable{source, d}.Subscribe
	}
}

type delayObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

type delayElement struct {
	Time  time.Time
	Value interface{}
}

func (obs delayObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Queue     queue.Queue
		Scheduled bool
		Completed bool
	}

	var doSchedule func(time.Duration)

	doSchedule = func(timeout time.Duration) {
		rx.Timer(timeout).Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}

			if critical.Enter(&x.Section) {
				x.Scheduled = false

				for x.Queue.Len() > 0 {
					if ctx.Err() != nil {
						break
					}

					t := x.Queue.Front().(delayElement)
					now := time.Now()
					if t.Time.After(now) {
						x.Scheduled = true
						doSchedule(t.Time.Sub(now))
						break
					}

					x.Queue.Pop()

					sink.Next(t.Value)
				}

				if x.Completed && x.Queue.Len() == 0 {
					sink.Complete()
				}

				critical.Leave(&x.Section)
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Queue.Push(
					delayElement{
						Time:  time.Now().Add(obs.Duration),
						Value: t.Value,
					},
				)

				if !x.Scheduled {
					x.Scheduled = true
					doSchedule(obs.Duration)
				}

				critical.Leave(&x.Section)

			case t.HasError:
				x.Queue.Init()

				critical.Close(&x.Section)

				sink(t)

			default:
				x.Completed = true

				if x.Queue.Len() > 0 {
					critical.Leave(&x.Section)
					break
				}

				critical.Close(&x.Section)

				sink(t)
			}
		}
	})
}
