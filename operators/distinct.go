package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Distinct creates an Observable that emits all items emitted by the source
// Observable that are distinct by comparison from previous items.
//
// If a keySelector function is provided, then it will project each value from
// the source Observable into a new value that it will check for equality with
// previously projected values. If a keySelector function is not provided, it
// will use each value from the source Observable directly with an equality
// check against previous values.
func Distinct() rx.Operator {
	return DistinctConfigure{}.Make()
}

// A DistinctConfigure is a configure for Distinct.
type DistinctConfigure struct {
	KeySelector func(interface{}) interface{}
}

// Make creates an Operator from this configure.
func (configure DistinctConfigure) Make() rx.Operator {
	if configure.KeySelector == nil {
		configure.KeySelector = func(val interface{}) interface{} { return val }
	}

	return func(source rx.Observable) rx.Observable {
		return distinctObservable{source, configure}.Subscribe
	}
}

type distinctObservable struct {
	Source rx.Observable
	DistinctConfigure
}

func (obs distinctObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	keys := make(map[interface{}]struct{})

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
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
