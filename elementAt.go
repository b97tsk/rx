package rx

import (
	"context"
)

type elementAtOperator struct {
	index           int
	defaultValue    interface{}
	hasDefaultValue bool
}

func (op elementAtOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	var (
		index           = op.index
		mutableObserver Observer
	)

	mutableObserver = func(t Notification) {
		switch {
		case t.HasValue:
			index--
			if index == -1 {
				mutableObserver = NopObserver
				t.Observe(ob)
				ob.Complete()
				cancel()
			}
		case t.HasError:
			t.Observe(ob)
			cancel()
		default:
			if op.hasDefaultValue {
				ob.Next(op.defaultValue)
				ob.Complete()
			} else {
				ob.Error(ErrOutOfRange)
			}
			cancel()
		}
	}

	source.Subscribe(ctx, func(t Notification) { t.Observe(mutableObserver) })

	return ctx, cancel
}

// ElementAt creates an Observable that emits the single value at the specified
// index in a sequence of emissions from the source Observable, if the
// specified index is out of range, notifies error ErrOutOfRange.
func (o Observable) ElementAt(index int) Observable {
	op := elementAtOperator{index: index}
	return o.Lift(op.Call)
}

// ElementAtOrDefault creates an Observable that emits the single value at the
// specified index in a sequence of emissions from the source Observable, if
// the specified index is out of range, emits the provided default value.
func (o Observable) ElementAtOrDefault(index int, defaultValue interface{}) Observable {
	op := elementAtOperator{
		index:           index,
		defaultValue:    defaultValue,
		hasDefaultValue: true,
	}
	return o.Lift(op.Call)
}
