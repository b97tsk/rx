package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type windowObservable struct {
	Source           rx.Observable
	WindowBoundaries rx.Observable
}

func (obs windowObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	type X struct {
		Window rx.Subject
	}
	window := rx.NewSubject()
	cx := make(chan *X, 1)
	cx <- &X{window}
	sink.Next(window.Observable)

	obs.WindowBoundaries.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Window.Complete()
				x.Window = rx.NewSubject()
				sink.Next(x.Window.Observable)
				cx <- x
			default:
				close(cx)
				t.Observe(x.Window.Observer)
				sink(t)
			}
		}
	})

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				t.Observe(x.Window.Observer)
				cx <- x
			default:
				close(cx)
				t.Observe(x.Window.Observer)
				sink(t)
			}
		}
	})
}

// Window branches out the source Observable values as a nested Observable
// whenever windowBoundaries emits.
//
// It's like Buffer, but emits a nested Observable instead of a slice.
func Window(windowBoundaries rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		obs := windowObservable{source, windowBoundaries}
		return rx.Create(obs.Subscribe)
	}
}
