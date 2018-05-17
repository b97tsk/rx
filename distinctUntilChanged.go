package rx

import (
	"context"
)

type distinctUntilChangedOperator struct {
	compare     func(interface{}, interface{}) bool
	keySelector func(interface{}) interface{}
}

func (op distinctUntilChangedOperator) Call(ctx context.Context, ob Observer, source Observable) (context.Context, context.CancelFunc) {
	var (
		key    interface{}
		hasKey bool
	)
	return source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			newKey := op.keySelector(t.Value)
			if hasKey && op.compare(key, newKey) {
				break
			}
			key = newKey
			hasKey = true
			t.Observe(ob)
		default:
			t.Observe(ob)
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
func (o Observable) DistinctUntilChanged() Observable {
	op := distinctUntilChangedOperator{defaultCompare, defaultKeySelector}
	return o.Lift(op.Call)
}
