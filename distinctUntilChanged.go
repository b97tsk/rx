package rx

import (
	"context"
)

type distinctUntilChangedOperator struct {
	Compare     func(interface{}, interface{}) bool
	KeySelector func(interface{}) interface{}
}

func (op distinctUntilChangedOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		key    interface{}
		hasKey bool
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			newKey := op.KeySelector(t.Value)
			if hasKey && op.Compare(key, newKey) {
				break
			}
			key = newKey
			hasKey = true
			sink(t)
		default:
			sink(t)
		}
	})
}

// DistinctUntilChanged creates an Observable that emits all items emitted by
// the source Observable that are distinct by comparison from the previous item.
//
// If a comparator function is provided, then it will be called for each item
// to test for whether or not that value should be emitted.
//
// If a comparator function is not provided, an equality check is used by default.
func (Operators) DistinctUntilChanged() OperatorFunc {
	return func(source Observable) Observable {
		op := distinctUntilChangedOperator{defaultCompare, defaultKeySelector}
		return source.Lift(op.Call)
	}
}
