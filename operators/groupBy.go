package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// GroupBy creates an Observable that groups the items emitted by the source
// Observable according to a specified criterion, and emits these grouped
// items as rx.GroupedObservables, one rx.GroupedObservable per group.
func GroupBy(keySelector func(interface{}) interface{}, groupFactory rx.DoubleFactory) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return groupByObservable{source, keySelector, groupFactory}.Subscribe
	}
}

type groupByObservable struct {
	Source       rx.Observable
	KeySelector  func(interface{}) interface{}
	GroupFactory rx.DoubleFactory
}

func (obs groupByObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	var groups = make(map[interface{}]rx.Observer)
	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			key := obs.KeySelector(t.Value)
			group, exists := groups[key]
			if !exists {
				d := obs.GroupFactory()
				group = d.Observer
				groups[key] = group
				sink.Next(rx.GroupedObservable{
					Observable: d.Observable,
					Key:        key,
				})
			}
			group.Sink(t)

		default:
			for _, group := range groups {
				group.Sink(t)
			}
			sink(t)
		}
	})
}
