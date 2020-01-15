package rx

import (
	"context"
)

// A DistinctConfigure is a configure for Distinct.
type DistinctConfigure struct {
	KeySelector func(interface{}) interface{}
}

// Use creates an Operator from this configure.
func (configure DistinctConfigure) Use() Operator {
	return func(source Observable) Observable {
		return distinctObservable{source, configure}.Subscribe
	}
}

type distinctObservable struct {
	Source Observable
	DistinctConfigure
}

func (obs distinctObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var keys = make(map[interface{}]struct{})
	return obs.Source.Subscribe(ctx, func(t Notification) {
		if t.HasValue {
			key := obs.KeySelector(t.Value)
			if _, exists := keys[key]; exists {
				return
			}
			keys[key] = struct{}{}
		}
		sink(t)
	})
}

// Distinct creates an Observable that emits all items emitted by the source
// Observable that are distinct by comparison from previous items.
//
// If a keySelector function is provided, then it will project each value from
// the source Observable into a new value that it will check for equality with
// previously projected values. If a keySelector function is not provided, it
// will use each value from the source Observable directly with an equality
// check against previous values.
func (Operators) Distinct() Operator {
	return DistinctConfigure{
		KeySelector: func(val interface{}) interface{} { return val },
	}.Use()
}
