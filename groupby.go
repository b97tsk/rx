package rx

// GroupBy groups the values emitted by the source Observable according to
// a specified criterion, and emits these grouped values as Pairs, one Pair
// per group.
func GroupBy[T any, K comparable](
	keySelector func(v T) K,
	groupFactory func() Subject[T],
) Operator[T, Pair[K, Observable[T]]] {
	return NewOperator(
		func(source Observable[T]) Observable[Pair[K, Observable[T]]] {
			return groupByObservable[T, K]{source, keySelector, groupFactory}.Subscribe
		},
	)
}

type groupByObservable[T any, K comparable] struct {
	source       Observable[T]
	keySelector  func(T) K
	groupFactory func() Subject[T]
}

func (ob groupByObservable[T, K]) Subscribe(c Context, o Observer[Pair[K, Observable[T]]]) {
	groups := make(map[K]Observer[T])

	ob.source.Subscribe(c, func(n Notification[T]) {
		switch n.Kind {
		case KindNext:
			key := ob.keySelector(n.Value)
			group, exists := groups[key]

			if !exists {
				g := ob.groupFactory()
				group = g.Observer
				groups[key] = group
				o.Next(NewPair(key, g.Observable))
			}

			group.Emit(n)

		case KindError, KindComplete:
			Try2(emitLastNotificationToGroups, groups, n, func() { o.Error(ErrOops) })

			switch n.Kind {
			case KindError:
				o.Error(n.Error)
			case KindComplete:
				o.Complete()
			}
		}
	})
}

func emitLastNotificationToGroups[T any, K comparable](groups map[K]Observer[T], n Notification[T]) {
	defer func() {
		if len(groups) != 0 {
			emitLastNotificationToGroups(groups, n)
		}
	}()

	for k, group := range groups {
		delete(groups, k)
		group.Emit(n)
	}
}
