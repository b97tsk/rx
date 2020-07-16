package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/subject"
)

type windowObservable struct {
	Source           rx.Observable
	WindowBoundaries rx.Observable
}

func (obs windowObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	type X struct {
		Window *subject.Subject
	}
	window := subject.NewSubject()
	cx := make(chan *X, 1)
	cx <- &X{window}
	sink.Next(window.Observable)

	obs.WindowBoundaries.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Window.Complete()
				x.Window = subject.NewSubject()
				sink.Next(x.Window.Observable)
				cx <- x
			default:
				close(cx)
				x.Window.Sink(t)
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
				x.Window.Sink(t)
				cx <- x
			default:
				close(cx)
				x.Window.Sink(t)
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
		return windowObservable{source, windowBoundaries}.Subscribe
	}
}
