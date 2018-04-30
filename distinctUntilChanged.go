package rx

import (
	"context"
)

type distinctUntilChangedOperator struct {
	source      Operator
	compare     func(interface{}, interface{}) bool
	keySelector func(interface{}) interface{}
}

func (op distinctUntilChangedOperator) ApplyOptions(options []Option) Operator {
	for _, opt := range options {
		switch t := opt.(type) {
		case compareOption:
			op.compare = t.Value
		case keySelectorOption:
			op.keySelector = t.Value
		default:
			panic(ErrUnsupportedOption)
		}
	}
	return op
}

func (op distinctUntilChangedOperator) Call(ctx context.Context, ob Observer) (context.Context, context.CancelFunc) {
	key := interface{}(nil)
	hasKey := false
	return op.source.Call(ctx, ObserverFunc(func(t Notification) {
		switch {
		case t.HasValue:
			newKey := op.keySelector(t.Value)
			if hasKey && op.compare(key, newKey) {
				break
			}
			key = newKey
			hasKey = true
			ob.Next(t.Value)
		case t.HasError:
			ob.Error(t.Value.(error))
		default:
			ob.Complete()
		}
	}))
}

// DistinctUntilChanged creates an Observable that emits all items emitted by
// the source Observable that are distinct by comparison from the previous item.
//
// If a comparator function is provided, then it will be called for each item
// to test for whether or not that value should be emitted.
//
// If a comparator function is not provided, an equality check is used by default.
func (o Observable) DistinctUntilChanged() Observable {
	op := distinctUntilChangedOperator{
		source:      o.Op,
		compare:     defaultCompare,
		keySelector: defaultKeySelector,
	}
	return Observable{op}
}
