package rx

import (
	"context"
)

type windowObservable struct {
	Source           Observable
	WindowBoundaries Observable
}

func (obs windowObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Window Subject
	}
	window := NewSubject()
	cx := make(chan *X, 1)
	cx <- &X{window}
	sink.Next(window.Observable)

	obs.WindowBoundaries.Subscribe(ctx, func(t Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Window.Complete()
				x.Window = NewSubject()
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
		return ctx, cancel
	}

	obs.Source.Subscribe(ctx, func(t Notification) {
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

	return ctx, cancel
}

// Window branches out the source Observable values as a nested Observable
// whenever windowBoundaries emits.
//
// It's like Buffer, but emits a nested Observable instead of a slice.
func (Operators) Window(windowBoundaries Observable) Operator {
	return func(source Observable) Observable {
		return windowObservable{source, windowBoundaries}.Subscribe
	}
}
