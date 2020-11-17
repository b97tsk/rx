package rx

import (
	"context"
)

// CongestingZip combines multiple Observables to create an Observable that
// emits the values of each of its input Observables as a slice.
//
// It's like Zip, but it congests subscribed Observables.
//
// For the purpose of allocation avoidance, slices emitted by the output
// Observable actually share the same underlying array.
func CongestingZip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return congestingZipObservable(observables).Subscribe
}

type congestingZipObservable []Observable

func (observables congestingZipObservable) Subscribe(ctx context.Context, sink Observer) {
	ctx, cancel := context.WithCancel(ctx)
	sink = sink.WithCancel(cancel)
	done := ctx.Done()

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
