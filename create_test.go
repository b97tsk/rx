package rx_test

import (
	"context"
	"testing"

	. "github.com/b97tsk/rx"
)

func TestCreate(t *testing.T) {
	obs := Create(
		func(ctx context.Context, sink Observer) {
			sink.Next("A")
			sink.Next("B")
			sink.Complete()
			sink.Next("C") // WARNING: buggy!
		},
	)
	subscribe(
		t,
		[]Observable{
			obs,
			obs.Pipe(operators.Mutex()), // Mutex gets rid of the extra "C".
		},
		"A", "B", Complete, "C",
		"A", "B", Complete,
	)
}
