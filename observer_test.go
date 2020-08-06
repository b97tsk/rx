package rx_test

import (
	"context"
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
	"github.com/b97tsk/rx/operators"
)

func TestObserver_ElementsOnly(t *testing.T) {
	Subscribe(
		t,
		rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				rx.Just("A", "B", "C").Subscribe(ctx, sink)
			},
		).Pipe(
			operators.Timeout(Step(1)),
		),
		"A", "B", "C", Completed,
	)
	Subscribe(
		t,
		rx.Observable(
			func(ctx context.Context, sink rx.Observer) {
				rx.Just("A", "B", "C").Subscribe(ctx, sink.ElementsOnly())
			},
		).Pipe(
			operators.Timeout(Step(1)),
		),
		"A", "B", "C", rx.ErrTimeout,
	)
}
