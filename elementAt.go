package rx

import (
	"context"
)

type elementAtOperator struct {
	source          Operator
	index           int
	defaultValue    interface{}
	hasDefaultValue bool
}

func (op elementAtOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case defaultValueOption:
			op.defaultValue = t.Value
			op.hasDefaultValue = true
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op elementAtOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	index := op.index

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			index--
			if index == -1 {
				mutable.Observer = NopObserver
				ob.Next(t.Value)
				ob.Complete()
				cancel()
			}
		case t.HasError:
			ob.Error(t.Value.(error))
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
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// ElementAt creates an Observable that emits the single value at the specified
// index in a sequence of emissions from the source Observable.
//
// ElementAt emits only the i-th value, then completes.
func (o Observable) ElementAt(index int) Observable {
	op := elementAtOperator{
		source: o.Op,
		index:  index,
	}
	return Observable{op}
}
