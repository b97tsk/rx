package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

type groupByObservable struct {
	Source         rx.Observable
	KeySelector    func(interface{}) interface{}
	SubjectFactory func() rx.Subject
}

func (obs groupByObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	var groups = make(map[interface{}]rx.Observer)
	return obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			key := obs.KeySelector(t.Value)
			group, exists := groups[key]
			if !exists {
				subject := obs.SubjectFactory()
				group = subject.Observer
				groups[key] = group
				sink.Next(rx.GroupedObservable{
					Observable: subject.Observable,
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

// GroupBy creates an Observable that groups the items emitted by the source
// Observable according to a specified criterion, and emits these grouped
// items as GroupedObservables, one GroupedObservable per group.
func GroupBy(keySelector func(interface{}) interface{}, subjectFactory func() rx.Subject) rx.Operator {
	return func(source rx.Observable) rx.Observable {
		return groupByObservable{source, keySelector, subjectFactory}.Subscribe
	}
}
