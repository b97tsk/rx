package rx

import (
	"context"
)

type windowWhenOperator struct {
	ClosingSelector func() Observable
}

func (op windowWhenOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	sink = Finally(sink, cancel)

	type X struct {
		Window Subject
	}
	window := NewSubject()
	cx := make(chan *X, 1)
	cx <- &X{window}
	sink.Next(window.Observable)

	var (
		openWindow     func()
		avoidRecursive avoidRecursiveCalls
	)

	openWindow = func() {
		if ctx.Err() != nil {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer Observer
		observer = func(t Notification) {
			observer = NopObserver
			cancel()
			if x, ok := <-cx; ok {
				if t.HasError {
					close(cx)
					t.Observe(x.Window.Observer)
					sink(t)
					return
				}
				x.Window.Complete()
				x.Window = NewSubject()
				sink.Next(x.Window.Observable)
				cx <- x
				avoidRecursive.Do(openWindow)
			}
		}

		closingNotifier := op.ClosingSelector()
		closingNotifier.Subscribe(ctx, observer.Notify)
	}

	avoidRecursive.Do(openWindow)

	if ctx.Err() != nil {
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

// WindowWhen branches out the source Observable values as a nested Observable
// using a factory function of closing Observables to determine when to start
// a new window.
//
// It's like BufferWhen, but emits a nested Observable instead of a slice.
func (Operators) WindowWhen(closingSelector func() Observable) OperatorFunc {
	return func(source Observable) Observable {
		op := windowWhenOperator{closingSelector}
		return source.Lift(op.Call)
	}
}
