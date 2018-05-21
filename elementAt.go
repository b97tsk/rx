package rx

import (
	"context"
)

type elementAtOperator struct {
	Index      int
	Default    interface{}
	HasDefault bool
}

func (op elementAtOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		index    = op.Index
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			index--
			if index == -1 {
				observer = NopObserver
				sink(t)
				sink.Complete()
				cancel()
			}
		case t.HasError:
			sink(t)
			cancel()
		default:
			if op.HasDefault {
				sink.Next(op.Default)
				sink.Complete()
			} else {
				sink.Error(ErrOutOfRange)
			}
			cancel()
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// ElementAt creates an Observable that emits the single value at the specified
// index in a sequence of emissions from the source Observable, if the
// specified index is out of range, notifies error ErrOutOfRange.
func (o Observable) ElementAt(index int) Observable {
	op := elementAtOperator{Index: index}
	return o.Lift(op.Call)
}

// ElementAtOrDefault creates an Observable that emits the single value at the
// specified index in a sequence of emissions from the source Observable, if
// the specified index is out of range, emits the provided default value.
func (o Observable) ElementAtOrDefault(index int, defaultValue interface{}) Observable {
	op := elementAtOperator{
		Index:      index,
		Default:    defaultValue,
		HasDefault: true,
	}
	return o.Lift(op.Call)
}
