package operators

import (
	"context"

	"github.com/b97tsk/rx"
)

// A ScanConfigure is a configure for Scan.
type ScanConfigure struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

// Use creates an Operator from this configure.
func (configure ScanConfigure) Use() rx.Operator {
	if configure.Accumulator == nil {
		panic("Scan: nil Accumulator")
	}
	return func(source rx.Observable) rx.Observable {
		return scanObservable{source, configure}.Subscribe
	}
}

type scanObservable struct {
	Source rx.Observable
	ScanConfigure
}

func (obs scanObservable) Subscribe(ctx context.Context, sink rx.Observer) (context.Context, context.CancelFunc) {
	sourceIndex := -1
	acc := struct {
		Value    interface{}
		HasValue bool
	}{obs.Seed, obs.HasSeed}
	return obs.Source.Subscribe(ctx, func(t rx.Notification) {
		switch {
		case t.HasValue:
			sourceIndex++

			if acc.HasValue {
				acc.Value = obs.Accumulator(acc.Value, t.Value, sourceIndex)
			} else {
				acc.Value = t.Value
				acc.HasValue = true
			}

			sink.Next(acc.Value)

		default:
			sink(t)
		}
	})
}

// Scan creates an Observable that applies an accumulator function over the
// source Observable, and emits each intermediate result, with an optional
// seed value.
//
// It's like Reduce, but emits the current accumulation whenever the source
// emits a value.
func Scan(accumulator func(interface{}, interface{}, int) interface{}) rx.Operator {
	return ScanConfigure{Accumulator: accumulator}.Use()
}
