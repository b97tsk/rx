package rx

import (
	"context"
)

type elementAtObservable struct {
	Source     Observable
	Index      int
	Default    interface{}
	HasDefault bool
}

func (obs elementAtObservable) Subscribe(ctx context.Context, sink Observer) {
	var (
		index    = obs.Index
		observer Observer
	)

	observer = func(t Notification) {
		switch {
		case t.HasValue:
			index--
			if index == -1 {
				observer = NopObserver
				sink(t)
				sink.Complete()
			}
		case t.HasError:
			sink(t)
		default:
			if obs.HasDefault {
				sink.Next(obs.Default)
				sink.Complete()
			} else {
				sink.Error(ErrOutOfRange)
			}
		}
	}

	obs.Source.Subscribe(ctx, observer.Notify)
}

// ElementAt creates an Observable that emits the single value at the specified
// index in a sequence of emissions from the source Observable, if the
// specified index is out of range, notifies error ErrOutOfRange.
func (Operators) ElementAt(index int) Operator {
	return func(source Observable) Observable {
		obs := elementAtObservable{
			Source: source,
			Index:  index,
		}
		return Create(obs.Subscribe)
	}
}

// ElementAtOrDefault creates an Observable that emits the single value at the
// specified index in a sequence of emissions from the source Observable, if
// the specified index is out of range, emits the provided default value.
func (Operators) ElementAtOrDefault(index int, defaultValue interface{}) Operator {
	return func(source Observable) Observable {
		obs := elementAtObservable{
			Source:     source,
			Index:      index,
			Default:    defaultValue,
			HasDefault: true,
		}
		return Create(obs.Subscribe)
	}
}
