package rx

import (
	"context"
)

type congestingZipObservable struct {
	Observables []Observable
}

func (obs congestingZipObservable) Subscribe(ctx context.Context, sink Observer) {
	done := ctx.Done()

	channels := make([]chan Notification, len(obs.Observables))

	for i := range channels {
		channels[i] = make(chan Notification)
	}

	go func() {
		for {
			nextValues := make([]interface{}, len(channels))
			for i := range channels {
				select {
				case <-done:
					return
				case t := <-channels[i]:
					switch {
					case t.HasValue:
						nextValues[i] = t.Value
					default:
						sink(t)
						return
					}
				}
			}
			sink.Next(nextValues)
		}
	}()

	for index, obs := range obs.Observables {
		c := channels[index]
		go obs.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case c <- t:
			}
		})
	}
}

// CongestingZip combines multiple Observables to create an Observable that
// emits the values of each of its input Observables as a slice.
//
// It's like Zip, but it congests subscribed Observables.
func CongestingZip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	obs := congestingZipObservable{observables}
	return Create(obs.Subscribe)
}
