package rx

// Contains emits a boolean to indicate whether any value of the source
// Observable satisfies a given condition.
func Contains[T any](cond func(v T) bool) Operator[T, bool] {
	if cond == nil {
		panic("cond == nil")
	}

	return contains(cond)
}

func contains[T any](cond func(v T) bool) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return containsObservable[T]{source, cond}.Subscribe
		},
	)
}

// ContainsElement emits a boolean to indicate whether the source Observable
// emits a given value.
func ContainsElement[T comparable](v T) Operator[T, bool] {
	return NewOperator(
		func(source Observable[T]) Observable[bool] {
			return containsElementObservable[T]{source, v}.Subscribe
		},
	)
}

type containsObservable[T any] struct {
	Source    Observable[T]
	Condition func(T) bool
}

func (obs containsObservable[T]) Subscribe(c Context, sink Observer[bool]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if obs.Condition(n.Value) {
				noop = true
				sink.Next(true)
				sink.Complete()
			}
		case KindError:
			sink.Error(n.Error)
		case KindComplete:
			sink.Next(false)
			sink.Complete()
		}
	})
}

type containsElementObservable[T comparable] struct {
	Source  Observable[T]
	Element T
}

func (obs containsElementObservable[T]) Subscribe(c Context, sink Observer[bool]) {
	c, cancel := c.WithCancel()
	sink = sink.OnLastNotification(cancel)

	var noop bool

	obs.Source.Subscribe(c, func(n Notification[T]) {
		if noop {
			return
		}

		switch n.Kind {
		case KindNext:
			if n.Value == obs.Element {
				noop = true
				sink.Next(true)
				sink.Complete()
			}
		case KindError:
			sink.Error(n.Error)
		case KindComplete:
			sink.Next(false)
			sink.Complete()
		}
	})
}
