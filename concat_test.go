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
		ErrComplete,
	).Case(
		rx.Concat(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
			rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
			rx.ConcatWith(
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Concat(
			func(_ rx.Context, sink rx.Observer[int]) {
				time.Sleep(Step(2))
				sink.Complete()
			},
			func(rx.Context, rx.Observer[int]) {
				panic("should not happen")
			},
		),
		context.DeadlineExceeded,
	)
}

func TestConcatMap(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.ConcatAll[rx.Observable[string]](),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMap(
				func(time.Time) rx.Observable[string] {
					return rx.Just("A")
				},
			),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMapTo[time.Time](rx.Just("A")),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.ConcatAll[rx.Observable[string]](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.NewObservable(
				func(_ rx.Context, sink rx.Observer[rx.Observable[string]]) {
					sink.Next(rx.Throw[string](ErrTest))
					sink.Complete()
				},
			),
			rx.ConcatAll[rx.Observable[string]](),
		),
		ErrTest,
	)
}

func TestConcatMapWithBuffering(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just(
				rx.Pipe1(rx.Just("A", "B"), AddLatencyToValues[string](3, 5)),
				rx.Pipe1(rx.Just("C", "D"), AddLatencyToValues[string](2, 4)),
				rx.Pipe1(rx.Just("E", "F"), AddLatencyToValues[string](1, 3)),
			),
			rx.ConcatAll[rx.Observable[string]]().WithBuffering(),
		),
		"A", "B", "C", "D", "E", "F", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMap(
				func(time.Time) rx.Observable[string] {
					return rx.Just("A")
				},
			).WithBuffering(),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMapTo[time.Time](rx.Just("A")).WithBuffering(),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Observable[string]](ErrTest),
			rx.ConcatAll[rx.Observable[string]]().WithBuffering(),
		),
		ErrTest,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeout(Step(1))
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe1(
			rx.Just[rx.Observable[string]](
				func(_ rx.Context, sink rx.Observer[string]) {
					time.Sleep(Step(2))
					sink.Complete()
				},
				func(rx.Context, rx.Observer[string]) {
					panic("should not happen")
				},
			),
			rx.ConcatAll[rx.Observable[string]]().WithBuffering(),
		),
		context.DeadlineExceeded,
	)
}
