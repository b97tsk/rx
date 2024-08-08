package rx_test

import (
	"testing"
	"time"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestOnBackpressure(t *testing.T) {
	t.Parallel()

	onNext := func(v int) {
		n := 4

		if v == 1 {
			// The first value arrives at odd step.
			// Decrease n to make sure follow-ups called at even steps so as to
			// pass the tests.
			n--
		}

		time.Sleep(Step(n))
	}

	t.Run("Buffer", func(t *testing.T) {
		t.Parallel()

		NewTestSuite[int](t).Case(
			rx.Pipe3(
				rx.Range(1, 8),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureBuffer[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, 5, 6, 7, ErrComplete,
		).Case(
			rx.Pipe3(
				rx.Range(1, 11),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureBuffer[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, rx.ErrBufferOverflow,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureBuffer[int](3),
			),
			ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Throw[int](ErrTest),
				rx.OnBackpressureBuffer[int](3),
			),
			ErrTest,
		).Case(
			rx.Pipe1(
				rx.Oops[int](ErrTest),
				rx.OnBackpressureBuffer[int](3),
			),
			rx.ErrOops, ErrTest,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureBuffer[int](0),
			),
			rx.ErrOops, "OnBackpressureBuffer: capacity < 1",
		)
	})

	t.Run("Congest", func(t *testing.T) {
		t.Parallel()

		NewTestSuite[int](t).Case(
			rx.Pipe3(
				rx.Range(1, 8),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureCongest[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, 5, 6, 7, ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureCongest[int](3),
			),
			ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Throw[int](ErrTest),
				rx.OnBackpressureCongest[int](3),
			),
			ErrTest,
		).Case(
			rx.Pipe1(
				rx.Oops[int](ErrTest),
				rx.OnBackpressureCongest[int](3),
			),
			rx.ErrOops, ErrTest,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureCongest[int](0),
			),
			rx.ErrOops, "OnBackpressureCongest: capacity < 1",
		)
	})

	t.Run("Drop", func(t *testing.T) {
		t.Parallel()

		NewTestSuite[int](t).Case(
			rx.Pipe3(
				rx.Range(1, 8),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureDrop[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, 5, 6, 7, ErrComplete,
		).Case(
			rx.Pipe3(
				rx.Range(1, 11),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureDrop[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, 5, 6, 7, 9, ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureDrop[int](3),
			),
			ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Throw[int](ErrTest),
				rx.OnBackpressureDrop[int](3),
			),
			ErrTest,
		).Case(
			rx.Pipe1(
				rx.Oops[int](ErrTest),
				rx.OnBackpressureDrop[int](3),
			),
			rx.ErrOops, ErrTest,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureDrop[int](0),
			),
			rx.ErrOops, "OnBackpressureDrop: capacity < 1",
		)
	})

	t.Run("Latest", func(t *testing.T) {
		t.Parallel()

		NewTestSuite[int](t).Case(
			rx.Pipe3(
				rx.Range(1, 8),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureLatest[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, 5, 6, 7, ErrComplete,
		).Case(
			rx.Pipe3(
				rx.Range(1, 11),
				AddLatencyToValues[int](1, 2),
				rx.OnBackpressureLatest[int](3),
				rx.DoOnNext(onNext),
			),
			1, 2, 3, 4, 6, 8, 9, 10, ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureLatest[int](3),
			),
			ErrComplete,
		).Case(
			rx.Pipe1(
				rx.Throw[int](ErrTest),
				rx.OnBackpressureLatest[int](3),
			),
			ErrTest,
		).Case(
			rx.Pipe1(
				rx.Oops[int](ErrTest),
				rx.OnBackpressureLatest[int](3),
			),
			rx.ErrOops, ErrTest,
		).Case(
			rx.Pipe1(
				rx.Empty[int](),
				rx.OnBackpressureLatest[int](0),
			),
			rx.ErrOops, "OnBackpressureLatest: capacity < 1",
		)
	})
}
