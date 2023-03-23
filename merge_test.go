package rx_test

import (
	"context"
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Merge[string](),
		ErrComplete,
	).Case(
		rx.Merge(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
			rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
		),
		"E", "C", "A", "F", "D", "B", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.MergeWith(
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
		),
		"E", "C", "A", "F", "D", "B", ErrComplete,
	)
}

func TestMerge2(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.MergeAll[rx.Observable[string]]().AsOperator(),
		),
		"E", "C", "A", "F", "D", "B", ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Range(0, 9),
			rx.MergeMap(
				func(v int) rx.Observable[int] {
					return rx.Pipe1(rx.Just(v), DelaySubscription[int](1))
				},
			).WithConcurrency(3).AsOperator(),
			rx.Reduce(0, func(v1, v2 int) int { return v1 + v2 }),
			ToString[int](),
		),
		"36", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MergeMapTo[time.Time](rx.Just("A")).AsOperator(),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[rx.Observable[string]](),
			rx.MergeAll[rx.Observable[string]]().AsOperator(),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.MergeAll[rx.Observable[string]]().AsOperator(),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			func(_ context.Context, sink rx.Observer[rx.Observable[string]]) {
				sink.Next(rx.Just("A", "B"))
				time.Sleep(Step(1))
				sink.Next(rx.Throw[string](ErrTest))
				time.Sleep(Step(1))
				sink.Next(rx.Just("C", "D"))
				time.Sleep(Step(1))
				sink.Next(rx.Just("E", "F"))
				sink.Complete()
			},
			rx.MergeAll[rx.Observable[string]]().AsOperator(),
		),
		"A", "B", ErrTest,
	)
}
