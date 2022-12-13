package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMap(t *testing.T) {
	t.Parallel()

	double := func(v int) int {
		return v * 2
	}

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[int](),
			rx.Map(double),
		),
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Range(1, 5),
			rx.Map(double),
		),
		2, 4, 6, 8, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Range(1, 5),
				rx.Throw[int](ErrTest),
			),
			rx.Map(double),
		),
		2, 4, 6, 8, ErrTest,
	)
}

func TestMapTo(t *testing.T) {
	t.Parallel()

	NewTestSuite[int](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.MapTo[string](42),
		),
		ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.MapTo[string](42),
		),
		42, 42, 42, ErrCompleted,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.MapTo[string](42),
		),
		42, 42, 42, ErrTest,
	)
}
