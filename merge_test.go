package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMerge(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Merge[string](),
		ErrCompleted,
	).Case(
		rx.Merge(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
			rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
		),
		"E", "C", "A", "F", "D", "B", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.MergeWith(
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
		),
		"E", "C", "A", "F", "D", "B", ErrCompleted,
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
		"E", "C", "A", "F", "D", "B", ErrCompleted,
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
		"36", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MergeMapTo[time.Time](rx.Just("A")).AsOperator(),
		),
		"A", ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Empty[rx.Observable[string]](),
			rx.MergeAll[rx.Observable[string]]().AsOperator(),
		),
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.MergeAll[rx.Observable[string]]().AsOperator(),
		),
		ErrTest,
	)
}
