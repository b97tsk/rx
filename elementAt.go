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

	sink = Finally(sink, cancel)

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
			}
		case t.HasError:
			sink(t)
		default:
			if op.HasDefault {
				sink.Next(op.Default)
				sink.Complete()
			} else {
				sink.Error(ErrOutOfRange)
			}
		}
	}

	source.Subscribe(ctx, observer.Notify)

	return ctx, cancel
}

// ElementAt creates an Observable that emits the single value at the specified
// index in a sequence of emissions from the source Observable, if the
// specified index is out of range, notifies error ErrOutOfRange.
func (Operators) ElementAt(index int) OperatorFunc {
	return func(source Observable) Observable {
		op := elementAtOperator{Index: index}
		return source.Lift(op.Call)
	}
}

// ElementAtOrDefault creates an Observable that emits the single value at the
// specified index in a sequence of emissions from the source Observable, if
// the specified index is out of range, emits the provided default value.
func (Operators) ElementAtOrDefault(index int, defaultValue interface{}) OperatorFunc {
	return func(source Observable) Observable {
		op := elementAtOperator{
			Index:      index,
			Default:    defaultValue,
			HasDefault: true,
		}
		return source.Lift(op.Call)
	}
}
