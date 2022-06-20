package rx

import (
	"context"
)

// GroupBy groups the items emitted by the source Observable according to
// a specified criterion, and emits these grouped items as Pairs, one Pair
// per group.
func GroupBy[T any, K comparable](
	keySelector func(v T) K,
	groupFactory func() Subject[T],
) Operator[T, Pair[K, Observable[T]]] {
	switch {
	case keySelector == nil:
		panic("keySelector == nil")
	case groupFactory == nil:
		panic("groupFactory == nil")
	}

	return groupBy(keySelector, groupFactory)
}

func groupBy[T any, K comparable](
	keySelector func(v T) K,
	groupFactory func() Subject[T],
) Operator[T, Pair[K, Observable[T]]] {
	return AsOperator(
		func(source Observable[T]) Observable[Pair[K, Observable[T]]] {
			return groupByObservable[T, K]{source, keySelector, groupFactory}.Subscribe
		},
	)
}

type groupByObservable[T any, K comparable] struct {
	Source       Observable[T]
	KeySelector  func(T) K
	GroupFactory func() Subject[T]
}

func (obs groupByObservable[T, K]) Subscribe(ctx context.Context, sink Observer[Pair[K, Observable[T]]]) {
	groups := make(map[K]Observer[T])

	obs.Source.Subscribe(ctx, func(n Notification[T]) {
		switch {
		case n.HasValue:
			key := obs.KeySelector(n.Value)

			group, exists := groups[key]

			if !exists {
				g := obs.GroupFactory()

				group = g.Observer
				groups[key] = group

				sink.Next(MakePair(key, g.Observable))
			}

			group.Sink(n)

		default:
			for _, group := range groups {
				group.Sink(n)
			}

			if n.HasError {
				sink.Error(n.Error)
			} else {
				sink.Complete()
			}
		}
	})
}
