package rx

import (
	"context"
)

type windowOperator struct {
	WindowBoundaries Observable
}

func (op windowOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Window Subject
	}
	window := NewSubject()
	cx := make(chan *X, 1)
	cx <- &X{window}
	sink.Next(window.Observable)

	op.WindowBoundaries.Subscribe(ctx, func(t Notification) {
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

	if isDone(ctx) {
		return Done()
	}

	source.Subscribe(ctx, func(t Notification) {
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
func (Operators) Window(windowBoundaries Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := windowOperator{windowBoundaries}
		return source.Lift(op.Call)
	}
}
