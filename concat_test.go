package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestConcat(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Concat[string](),
		ErrCompleted,
	).Case(
		rx.Concat(
			rx.Pipe(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.Pipe(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
			rx.Pipe(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
		),
		"A", "B", "C", "D", "E", "F", ErrCompleted,
	).Case(
		rx.Pipe(
			rx.Pipe(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.ConcatWith(
				rx.Pipe(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
		),
		"A", "B", "C", "D", "E", "F", ErrCompleted,
	)

	ctx, cancel := context.WithTimeout(context.Background(), Step(1))
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Concat(
			func(_ context.Context, sink rx.Observer[int]) {
				time.Sleep(Step(2))
				sink.Complete()
			},
			func(context.Context, rx.Observer[int]) {
				panic("should not happen")
			},
		),
		context.DeadlineExceeded,
	)
}
