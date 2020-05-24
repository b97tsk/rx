package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/x/queue"
)

type delayObservable struct {
	Source   rx.Observable
	Duration time.Duration
}

type delayValue struct {
	Time  time.Time
	Value interface{}
}

func (obs delayObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Queue           queue.Queue
		Scheduled       bool
		SourceCompleted bool
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var doSchedule func(time.Duration)

	doSchedule = func(timeout time.Duration) {
		rx.Timer(timeout).Subscribe(ctx, func(t rx.Notification) {
			if t.HasValue {
				return
			}
			if x, ok := <-cx; ok {
				x.Scheduled = false
				for x.Queue.Len() > 0 {
					if ctx.Err() != nil {
						break
					}
					t := x.Queue.Front().(delayValue)
					now := time.Now()
					if t.Time.After(now) {
						x.Scheduled = true
						doSchedule(t.Time.Sub(now))
						break
					}
					x.Queue.PopFront()
					sink.Next(t.Value)
				}
				if x.SourceCompleted && x.Queue.Len() == 0 {
					sink.Complete()
				}
				cx <- x
			}
		})
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Queue.PushBack(
					delayValue{
						Time:  time.Now().Add(obs.Duration),
						Value: t.Value,
					},
				)
				if !x.Scheduled {
					x.Scheduled = true
					doSchedule(obs.Duration)
				}
				cx <- x
			case t.HasError:
				x.Queue.Init()
				close(cx)
				sink(t)
			default:
				x.SourceCompleted = true
				if x.Queue.Len() > 0 {
					cx <- x
					break
				}
				close(cx)
				sink(t)
			}
		}
	})
}

// Delay delays the emission of items from the source Observable by a given
// timeout.
func Delay(timeout time.Duration) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := delayObservable{source, timeout}
		return rx.Create(obs.Subscribe)
	}
}
