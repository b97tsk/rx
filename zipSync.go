package rx

import (
	"context"
)

// ZipSync combines multiple Observables to create an Observable that emits
// the values of each of its input Observables as a slice.
//
// It's like Zip, but it does not buffer subscribed Observables.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func ZipSync(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}

	return zipSyncObservable(observables).Subscribe
}

type zipSyncObservable []Observable

func (observables zipSyncObservable) Subscribe(ctx context.Context, sink Observer) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = sink.WithCancel(cancel)

	channels := make([]chan Notification, len(observables))

	for i := range channels {
		channels[i] = make(chan Notification)
	}

	go func() {
		values := make([]interface{}, len(channels))

		for {
			for i := range channels {
				select {
				case <-done:
					return
				case t := <-channels[i]:
					switch {
					case t.HasValue:
						values[i] = t.Value
					default:
						sink(t)
						return
					}
				}
			}

			sink.Next(values)
		}
	}()

	for i, obs := range observables {
		c := channels[i]

		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case c <- t:
			}
		})
	}
}
