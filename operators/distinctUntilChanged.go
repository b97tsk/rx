package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// DistinctUntilChanged emits all items emitted by the source that are
// distinct by comparison from the previous item.
//
// If a comparator function is provided, then it will be called for each
// item to test for whether or not that value should be emitted.
//
// If a comparator function is not provided, an equality check is used by
// default.
func DistinctUntilChanged() rx.Operator {
	return DistinctUntilChangedConfig{}.Make()
}

// A DistinctUntilChangedConfig is a configuration for DistinctUntilChanged.
type DistinctUntilChangedConfig struct {
	Compare     func(interface{}, interface{}) bool
	KeySelector func(interface{}) interface{}
}

// Make creates an Operator from this configuration.
func (config DistinctUntilChangedConfig) Make() rx.Operator {
	if config.Compare == nil {
		config.Compare = func(v1, v2 interface{}) bool { return v1 == v2 }
	}

	if config.KeySelector == nil {
		config.KeySelector = func(val interface{}) interface{} { return val }
	}

	return func(source rx.Observable) rx.Observable {
		return distinctUntilChangedObservable{source, config}.Subscribe
	}
}

type distinctUntilChangedObservable struct {
	Source rx.Observable
	DistinctUntilChangedConfig
}

func (obs distinctUntilChangedObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var key struct {
		Value    interface{}
		HasValue bool
	}

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		if t.HasValue {
			keyValue := obs.KeySelector(t.Value)

			if key.HasValue && obs.Compare(key.Value, keyValue) {
				return
			}

			key.Value = keyValue
			key.HasValue = true
		}

		sink(t)
	})
}
