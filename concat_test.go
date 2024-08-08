package rx_test

import (
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

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	NewTestSuite[int](t).WithContext(ctx).Case(
		rx.Concat(
			func(_ rx.Context, o rx.Observer[int]) {
				time.Sleep(Step(2))
				o.Complete()
			},
			rx.Oops[int]("should not happen"),
		),
		ErrTest,
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
			rx.Oops[rx.Observable[string]](ErrTest),
			rx.ConcatAll[rx.Observable[string]](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.NewObservable(
				func(_ rx.Context, o rx.Observer[rx.Observable[string]]) {
					o.Next(rx.Throw[string](ErrTest))
					o.Complete()
				},
			),
			rx.ConcatAll[rx.Observable[string]](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMap(func(time.Time) rx.Observable[string] { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
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
	).Case(
		rx.Pipe1(
			rx.Oops[rx.Observable[string]](ErrTest),
			rx.ConcatAll[rx.Observable[string]]().WithBuffering(),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Timer(Step(1)),
			rx.ConcatMap(func(time.Time) rx.Observable[string] { panic(ErrTest) }).WithBuffering(),
		),
		rx.ErrOops, ErrTest,
	)

	ctx, cancel := rx.NewBackgroundContext().WithTimeoutCause(Step(1), ErrTest)
	defer cancel()

	NewTestSuite[string](t).WithContext(ctx).Case(
		rx.Pipe1(
			rx.Just[rx.Observable[string]](
				func(_ rx.Context, o rx.Observer[string]) {
					time.Sleep(Step(2))
					o.Complete()
				},
				rx.Oops[string]("should not happen"),
			),
			rx.ConcatAll[rx.Observable[string]]().WithBuffering(),
		),
		ErrTest,
	)
}
