package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestMaterialize(t *testing.T) {
	t.Parallel()

	count := func(n int, _ rx.Notification[string]) int {
		return n + 1
	}

	NewTestSuite[int](t).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		1, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		4, ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Materialize[string](),
			rx.Reduce(0, count),
		),
		4, ErrComplete,
	).Case(
		rx.Pipe3(
			rx.Empty[string](),
			rx.Materialize[string](),
			rx.OnNext(func(rx.Notification[string]) { panic(ErrTest) }),
			rx.Reduce(0, count),
		),
		rx.ErrOops, ErrTest,
	)
}
