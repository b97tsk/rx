package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// SkipLast skip the last count values emitted by the source.
func SkipLast(count int) rx.Operator {
	if count <= 0 {
		return noop
	}

	return func(source rx.Observable) rx.Observable {
		return skipLastObservable{source, count}.Subscribe
	}
}

type skipLastObservable struct {
	Source rx.Observable
	Count  int
}

func (obs skipLastObservable) Subscribe(ctx context.Context, sink rx.Observer) {
	buffer := make([]interface{}, obs.Count)
	bufferSize := obs.Count

	var index, count int

	obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			if count < bufferSize {
				count++
			} else {
				sink.Next(buffer[index])
			}

			buffer[index] = t.Value

			index = (index + 1) % bufferSize

		default:
			sink(t)
		}
	})
}
