package operators

import (
	"context"
	"time"

	"github.com/b97tsk/rx"
)

type debounceObservable struct {
	Source           rx.Observable
	DurationSelector func(interface{}) rx.Observable
}

func (obs debounceObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Latest struct {
			Value    interface{}
			HasValue bool
		}
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var (
		scheduleCtx    context.Context
		scheduleCancel context.CancelFunc
	)

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Latest.Value = t.Value
				x.Latest.HasValue = true

				cx <- x

				if scheduleCancel != nil {
					scheduleCancel()
				}

				scheduleCtx, scheduleCancel = context.WithCancel(ctx)

				var observer rx.Observer
				observer = func(t rx.Notification) {
					observer = rx.Noop
					scheduleCancel()
					if x, ok := <-cx; ok {
						if t.HasError {
							close(cx)
							sink(t)
							return
						}
						if x.Latest.HasValue {
							sink.Next(x.Latest.Value)
							x.Latest.HasValue = false
						}
						cx <- x
					}
				}

				obs := obs.DurationSelector(t.Value)
				obs.Subscribe(scheduleCtx, observer.Sink)

			case t.HasError:
				close(cx)
				sink(t)

			default:
				close(cx)
				if x.Latest.HasValue {
					sink.Next(x.Latest.Value)
				}
				sink(t)
			}
		}
	})
}

// Debounce creates an Observable that emits a value from the source Observable
// only after a particular time span, determined by another Observable, has
// passed without another source emission.
//
// It's like DebounceTime, but the time span of emission silence is determined
// by a second Observable.
func Debounce(durationSelector func(interface{}) rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := debounceObservable{source, durationSelector}
		return rx.Create(obs.Subscribe)
	}
}

// DebounceTime creates an Observable that emits a value from the source
// Observable only after a particular time span has passed without another
// source emission.
func DebounceTime(d time.Duration) rx.Operator {
	obsTimer := rx.Timer(d)
	durationSelector := func(interface{}) rx.Observable { return obsTimer }
	return Debounce(durationSelector)
}
