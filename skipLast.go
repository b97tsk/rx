package rx

import (
	"context"
)

type skipLastObservable struct {
	Source Observable
	Count  int
}

func (obs skipLastObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		buffer     = make([]interface{}, obs.Count)
		bufferSize = obs.Count
		index      int
		count      int
	)
	return obs.Source.Subscribe(ctx, func(t Notification) {
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

// SkipLast creates an Observable that skip the last count values emitted by
// the source Observable.
func (Operators) SkipLast(count int) Operator {
	return func(source Observable) Observable {
		if count <= 0 {
			return source
		}
		return skipLastObservable{source, count}.Subscribe
	}
}
