package operators

import (
	"context"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/internal/critical"
	"github.com/b97tsk/rx/internal/norec"
)

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

type bufferWhenObservable struct {
	Source          rx.Observable
	ClosingSelector func() rx.Observable
}

func (obs bufferWhenObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	ctx, cancel := context.WithCancel(ctx)

	sink = sink.WithCancel(cancel)

	var x struct {
		critical.Section
		Buffer []interface{}
	}

	var openBuffer func()

	openBuffer = norec.Wrap(func() {
		if ctx.Err() != nil {
			return
		}

		ctx, cancel := context.WithCancel(ctx)

		var observer rx.Observer

		observer = func(t rx.Notification) {
			observer = rx.Noop

			cancel()

			if critical.Enter(&x.Section) {
				switch {
				case t.HasValue:
					buffer := x.Buffer
					x.Buffer = nil
					sink.Next(buffer)

					critical.Leave(&x.Section)

					openBuffer()

				case t.HasError:
					critical.Close(&x.Section)

					sink(t)

				default:
					critical.Leave(&x.Section)
				}
			}
		}

		closingNotifier := obs.ClosingSelector()

		closingNotifier.Subscribe(ctx, observer.Sink)
	})

	openBuffer()

	if ctx.Err() != nil {
		return
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if critical.Enter(&x.Section) {
			switch {
			case t.HasValue:
				x.Buffer = append(x.Buffer, t.Value)

				critical.Leave(&x.Section)

			case t.HasError:
				critical.Close(&x.Section)

				sink(t)

			default:
				critical.Close(&x.Section)

				sink.Next(x.Buffer)
				sink(t)
			}
		}
	})
}
