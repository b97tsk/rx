package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestCreate(t *testing.T) {
	obs := rx.Create(
		func(ctx context.Context, sink rx.Observer) {
			sink.Next("A")
			sink.Next("B")
			sink.Complete()
			sink.Next("C") // WARNING: buggy!
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			obs,
			obs.Pipe(operators.Mutex()), // Mutex gets rid of the extra "C".
		},
		[][]interface{}{
			{"A", "B", rx.Complete, "C"},
			{"A", "B", rx.Complete},
		},
	)
}
