package rx

import (
	"context"
)

type congestingZipObservable struct {
	Observables []Observable
}

func (obs congestingZipObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	sink = Finally(sink, cancel)

	length := len(obs.Observables)
	channels := make([]chan Notification, length)

	for i := 0; i < length; i++ {
		channels[i] = make(chan Notification)
	}

	go func() {
		for {
			nextValues := make([]interface{}, length)

			for i := 0; i < length; i++ {
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

	return ctx, cancel
}

// CongestingZip combines multiple Observables to create an Observable that
// emits the values of each of its input Observables as a slice.
//
// It's like Zip, but it congests subscribed Observables.
func CongestingZip(observables ...Observable) Observable {
	if len(observables) == 0 {
		return Empty()
	}
	return congestingZipObservable{observables}.Subscribe
}

// CongestingZipAll converts a higher-order Observable into a first-order
// Observable by waiting for the outer Observable to complete, then applying
// CongestingZip.
//
// CongestingZipAll flattens an Observable-of-Observables by applying
// CongestingZip when the Observable-of-Observables completes.
//
// It's like ZipAll, but it congests subscribed Observables.
func (Operators) CongestingZipAll() Operator {
	return ToObservablesConfigure{CongestingZip}.Use()
}
