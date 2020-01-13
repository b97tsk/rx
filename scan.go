package rx

import (
	"context"
)

// A ScanConfigure is a configure for Scan.
type ScanConfigure struct {
	Accumulator func(interface{}, interface{}, int) interface{}
	Seed        interface{}
	HasSeed     bool
}

// Use creates an Operator from this configure.
func (configure ScanConfigure) Use() Operator {
	return func(source Observable) Observable {
		return scanObservable{source, configure}.Subscribe
	}
}

type scanObservable struct {
	Source Observable
	ScanConfigure
}

func (obs scanObservable) Subscribe(ctx context.Context, sink Observer) (context.Context, context.CancelFunc) {
	var (
		seed       = obs.Seed
		hasSeed    = obs.HasSeed
		outerIndex = -1
	)
	return obs.Source.Subscribe(ctx, func(t Notification) {
		switch {
		case t.HasValue:
			outerIndex++

			if hasSeed {
				seed = obs.Accumulator(seed, t.Value, outerIndex)
			} else {
				seed = t.Value
				hasSeed = true
			}

			sink.Next(seed)

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
func (Operators) Scan(accumulator func(interface{}, interface{}, int) interface{}) Operator {
	return ScanConfigure{Accumulator: accumulator}.Use()
}
