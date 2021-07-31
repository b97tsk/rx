package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// Distinct emits all items emitted by the source that are distinct by
// comparison from previous items.
//
// If a keySelector function is provided, then it will project each value
// from the source into a new value that it will check for equality with
// previously projected values. If a keySelector function is not provided,
// it will use each value from the source directly with an equality check
// against previous values.
func Distinct() rx.Operator {
	return DistinctConfig{}.Make()
}

// A DistinctConfig is a configuration for Distinct.
type DistinctConfig struct {
	KeySelector func(interface{}) interface{}
}

// Make creates an Operator from this configuration.
func (config DistinctConfig) Make() rx.Operator {
	if config.KeySelector == nil {
		config.KeySelector = func(val interface{}) interface{} { return val }
	}

	return func(source rx.Observable) rx.Observable {
		return distinctObservable{source, config}.Subscribe
	}
}

type distinctObservable struct {
	Source rx.Observable
	DistinctConfig
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
