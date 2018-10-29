package rx

import (
	"context"
)

type congestingZipOperator struct {
	Observables []Observable
}

func (op congestingZipOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	done := ctx.Done()

	length := len(op.Observables)
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
						cancel()
						return
					}
				}
			}

			sink.Next(nextValues)
		}
	}()

	for index, obsv := range op.Observables {
		c := channels[index]
		go obsv.Subscribe(ctx, func(t Notification) {
			select {
			case <-done:
			case c <- t:
			}
		})
	}

	return ctx, cancel
}

type congestingZipAllOperator struct{}

func (op congestingZipAllOperator) Call(ctx context.Context, sink Observer, source Observable) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	toObservablesOperator(op).Call(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			observables := t.Value.([]Observable)

			if len(observables) == 0 {
				sink.Complete()
				cancel()
				break
			}

			zip := congestingZipOperator{observables}
			zip.Call(ctx, Finally(sink, cancel), Observable{})

		case t.HasError:
			sink(t)
			cancel()

		default:
		}
	}, source)

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
	op := congestingZipOperator{observables}
	return Observable{}.Lift(op.Call)
}

// CongestingZipAll converts a higher-order Observable into a first-order
// Observable by waiting for the outer Observable to complete, then applying
// CongestingZip.
//
// CongestingZipAll flattens an Observable-of-Observables by applying
// CongestingZip when the Observable-of-Observables completes.
//
// It's like ZipAll, but it congests subscribed Observables.
func (Operators) CongestingZipAll() OperatorFunc {
	return func(source Observable) Observable {
		op := congestingZipAllOperator{}
		return source.Lift(op.Call)
	}
}
