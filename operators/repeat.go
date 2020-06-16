package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/misc"
)

type repeatObservable struct {
	Source rx.Observable
	Count  int
}

func (obs repeatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var (
		count          = obs.Count
		observer       rx.Observer
		avoidRecursion misc.AvoidRecursion
	)

	subscribe := func() {
		obs.Source.Subscribe(ctx, observer)
	}

	observer = func(t rx.Notification) {
		if t.HasValue || t.HasError || count == 0 {
			sink(t)
			return
		}
		if count > 0 {
			count--
		}
		avoidRecursion.Do(subscribe)
	}

	avoidRecursion.Do(subscribe)
}

// Repeat creates an Observable that repeats the stream of items emitted by the
// source Observable at most count times.
func Repeat(count int) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		if count == 0 {
			return rx.Empty()
		}
		if count == 1 {
			return source
		}
		if count > 0 {
			count--
		}
		obs := repeatObservable{source, count}
		return rx.Create(obs.Subscribe)
	}
}

// RepeatForever creates an Observable that repeats the stream of items emitted
// by the source Observable forever.
func RepeatForever() rx.Operator {
	return Repeat(-1)
}
