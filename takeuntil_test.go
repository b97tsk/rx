package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestTakeUntil(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Just(42)),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Empty[int]()),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Never[int]()),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](rx.Throw[int](ErrTest)),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](
				rx.Pipe1(
					rx.Just(42),
					DelaySubscription[int](3),
				),
			),
		),
		"A", "B", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](
				rx.Pipe1(
					rx.Empty[int](),
					DelaySubscription[int](3),
				),
			),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			AddLatencyToValues[string](0, 2),
			rx.TakeUntil[string](
				rx.Pipe1(
					rx.Throw[int](ErrTest),
					DelaySubscription[int](3),
				),
			),
		),
		"A", "B", ErrTest,
	)

	t.Run("Oops", func(t *testing.T) {
		defer func() {
			NewTestSuite[string](t).Case(rx.Oops[string](recover()), rx.ErrOops, ErrTest)
		}()
		rx.Pipe1(
			rx.Empty[string](),
			rx.TakeUntil[string](
				func(_ rx.Context, sink rx.Observer[int]) {
					defer sink.Complete()
					panic(ErrTest)
				},
			),
		).Subscribe(rx.NewBackgroundContext(), rx.Noop[string])
	})
}
