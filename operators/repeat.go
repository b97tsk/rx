package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/norec"
)

// Repeat repeats the stream of items emitted by the source at most count
// times.
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

		return repeatObservable{source, count}.Subscribe
	}
}

// RepeatForever repeats the stream of items emitted by the source forever.
func RepeatForever() rx.Operator {
	return Repeat(-1)
}

type repeatObservable struct {
	Source rx.Observable
	Count  int
}

func (obs repeatObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var observer rx.Observer

	subscribeToSource := norec.Wrap(func() {
		obs.Source.Subscribe(ctx, observer)
	})

	count := obs.Count

	observer = func(t rx.Notification) {
		if t.HasValue || t.HasError || count == 0 {
			sink(t)
			return
		}

		if count > 0 {
			count--
		}

		subscribeToSource()
	}

	subscribeToSource()
}
