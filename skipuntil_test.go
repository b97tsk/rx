package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSkipUntil(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Just(42)),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Empty[int]()),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Never[int]()),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Throw[int](ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](rx.Oops[int](ErrTest)),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe1(
					rx.Just(42),
					DelaySubscription[int](3),
				),
			),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe1(
					rx.Empty[int](),
					DelaySubscription[int](3),
				),
			),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe1(
					rx.Throw[int](ErrTest),
					DelaySubscription[int](3),
				),
			),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.SkipUntil[string](
				rx.Pipe1(
					rx.Oops[int](ErrTest),
					DelaySubscription[int](3),
				),
			),
		),
		rx.ErrOops, ErrTest,
	)

	t.Run("Oops", func(t *testing.T) {
		defer func() {
			NewTestSuite[string](t).Case(rx.Oops[string](recover()), rx.ErrOops, ErrTest)
		}()
		rx.Pipe1(
			rx.Empty[string](),
			rx.SkipUntil[string](
				func(_ rx.Context, o rx.Observer[int]) {
					defer o.Complete()
					panic(ErrTest)
				},
			),
		).Subscribe(rx.NewBackgroundContext(), rx.Noop[string])
	})
}
