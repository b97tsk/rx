package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/misc"
)

type bufferWhenObservable struct {
	Source          rx.Observable
	ClosingSelector func() rx.Observable
}

func (obs bufferWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)

	type X struct {
		Buffers []interface{}
	}
	cx := make(chan *X, 1)
	cx <- &X{}

	var (
		openBuffer     func()
		avoidRecursion misc.AvoidRecursion
	)

	openBuffer = func() {
		if ctx.Err() != nil {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer
		observer = func(t rx.Notification) {
			observer = rx.Noop
			cancel()
			if x, ok := <-cx; ok {
				if t.HasError {
					close(cx)
					sink(t)
					return
				}
				sink.Next(x.Buffers)
				x.Buffers = nil
				cx <- x
				avoidRecursion.Do(openBuffer)
			}
		}

		closingNotifier := obs.ClosingSelector()
		closingNotifier.Subscribe(ctx, observer.Sink)
	}

	avoidRecursion.Do(openBuffer)

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if x, ok := <-cx; ok {
			switch {
			case t.HasValue:
				x.Buffers = append(x.Buffers, t.Value)
				cx <- x
			default:
				close(cx)
				sink.Next(x.Buffers)
				sink(t)
			}
		}
	})
}

// BufferWhen buffers the source Observable values, using a factory function
// of closing Observables to determine when to close, emit, and reset the
// buffer.
//
// BufferWhen collects values from the past as a slice. When it starts
// collecting values, it calls a function that returns an Observable that
// tells when to close the buffer and restart collecting.
func BufferWhen(closingSelector func() rx.Observable) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return bufferWhenObservable{source, closingSelector}.Subscribe
	}
}
