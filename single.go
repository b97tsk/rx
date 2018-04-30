package rx

import (
	"context"
)

type singleOperator struct {
	source    Operator
	predicate func(interface{}, int) bool
}

func (op singleOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case predicateOption:
			op.predicate = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op singleOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	value := interface{}(nil)
	hasValue := false
	outerIndex := -1

	mutable := MutableObserver{}

	mutable.Observer = ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if !op.predicate(t.Value, outerIndex) {
				break
			}

			if hasValue {
				mutable.Observer = NopObserver
				ob.Error(ErrNotSingle)
				cancel()
			} else {
				value = t.Value
				hasValue = true
			}

		case t.HasError:
			ob.Error(t.Value.(error))
			cancel()

		default:
			if hasValue {
				ob.Next(value)
				ob.Complete()
			} else {
				ob.Error(ErrEmpty)
			}
			cancel()
		}
	})

	op.source.Call(ctx, &mutable)

	return ctx, cancel
}

// Single creates an Observable that emits the single item emitted by the
// source Observable that matches a specified predicate, if that Observable
// emits one such item. If the source Observable emits more than one such item
// or no such items, notify of an ErrNotSingle or ErrEmpty respectively.
func (o Observable) Single() Observable {
	op := singleOperator{
		source:    o.Op,
		predicate: defaultPredicate,
	}
	return Observable{op}
}
