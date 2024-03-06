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
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.MergeWith(rx.Just("A")),
		),
		ErrTest,
	)
}

func TestMergeMap(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.MergeAll[rx.Observable[string]](),
		),
		"E", "C", "A", "F", "D", "B", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(
				rx.Just("A", "B"),
				rx.Just("C", "D"),
				rx.Just("E", "F"),
			),
			rx.MergeAll[rx.Observable[string]]().WithPassiveGo(),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Range(0, 9),
			rx.MergeMap(
				func(v int) rx.Observable[int] {
					return rx.Pipe1(rx.Just(v), DelaySubscription[int](1))
				},
			).WithConcurrency(3),
			rx.Reduce(0, func(v1, v2 int) int { return v1 + v2 }),
			ToString[int](),
		),
		"36", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MergeMapTo[time.Time](rx.Just("A")),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[rx.Observable[string]](),
			rx.MergeAll[rx.Observable[string]](),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.MergeAll[rx.Observable[string]](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			func(_ rx.Context, sink rx.Observer[rx.Observable[string]]) {
				sink.Next(rx.Just("A", "B"))
				time.Sleep(Step(1))
				sink.Next(rx.Throw[string](ErrTest))
				time.Sleep(Step(1))
				sink.Next(rx.Just("C", "D"))
				time.Sleep(Step(1))
				sink.Next(rx.Just("E", "F"))
				sink.Complete()
			},
			rx.MergeAll[rx.Observable[string]](),
		),
		"A", "B", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[rx.Observable[string]](),
			rx.MergeAll[rx.Observable[string]]().WithConcurrency(0),
		),
		rx.ErrOops, "MergeMap: Concurrency == 0",
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MergeMap(func(time.Time) rx.Observable[string] { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestMergeMapWithBuffering(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.MergeAll[rx.Observable[string]]().WithBuffering(),
		),
		"E", "C", "A", "F", "D", "B", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just(
				rx.Just("A", "B"),
				rx.Just("C", "D"),
				rx.Just("E", "F"),
			),
			rx.MergeAll[rx.Observable[string]]().WithBuffering().WithPassiveGo(),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Range(0, 9),
			rx.MergeMap(
				func(v int) rx.Observable[int] {
					return rx.Pipe1(rx.Just(v), DelaySubscription[int](1))
				},
			).WithBuffering().WithConcurrency(3),
			rx.Reduce(0, func(v1, v2 int) int { return v1 + v2 }),
			ToString[int](),
		),
		"36", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MergeMapTo[time.Time](rx.Just("A")).WithBuffering(),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[rx.Observable[string]](),
			rx.MergeAll[rx.Observable[string]]().WithBuffering(),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.MergeAll[rx.Observable[string]]().WithBuffering(),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			func(_ rx.Context, sink rx.Observer[rx.Observable[string]]) {
				sink.Next(rx.Just("A", "B"))
				time.Sleep(Step(1))
				sink.Next(rx.Throw[string](ErrTest))
				time.Sleep(Step(1))
				sink.Next(rx.Just("C", "D"))
				time.Sleep(Step(1))
				sink.Next(rx.Just("E", "F"))
				sink.Complete()
			},
			rx.MergeAll[rx.Observable[string]]().WithBuffering(),
		),
		"A", "B", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.MergeMap(func(time.Time) rx.Observable[string] { panic(ErrTest) }).WithBuffering(),
		),
		rx.ErrOops, ErrTest,
	)
}
