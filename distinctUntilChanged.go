package rx

import (
	"context"
)

// A DistinctUntilChangedConfigure is a configure for DistinctUntilChanged.
type DistinctUntilChangedConfigure struct {
	Compare     func(interface{}, interface{}) bool
	KeySelector func(interface{}) interface{}
}

// Use creates an Operator from this configure.
func (configure DistinctUntilChangedConfigure) Use() Operator {
	return func(source Observable) Observable {
		return distinctUntilChangedObservable{source, configure}.Subscribe
	}
}

type distinctUntilChangedObservable struct {
	Source Observable
	DistinctUntilChangedConfigure
}

func (obs distinctUntilChangedObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		key    interface{}
		hasKey bool
	)
	return obs.Source.Subscribe(ctx, func(t Notification) {
		if t.HasValue {
			newKey := obs.KeySelector(t.Value)
			if hasKey && obs.Compare(key, newKey) {
				return
			}
			key = newKey
			hasKey = true
		}
		sink(t)
	})
}

// DistinctUntilChanged creates an Observable that emits all items emitted by
// the source Observable that are distinct by comparison from the previous item.
//
// If a comparator function is provided, then it will be called for each item
// to test for whether or not that value should be emitted.
//
// If a comparator function is not provided, an equality check is used by default.
func (Operators) DistinctUntilChanged() Operator {
	return DistinctUntilChangedConfigure{
		Compare:     func(v1, v2 interface{}) bool { return v1 == v2 },
		KeySelector: func(val interface{}) interface{} { return val },
	}.Use()
}
