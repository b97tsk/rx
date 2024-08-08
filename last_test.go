package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestLast(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.Last[string](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.Last[string](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.Last[string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Last[string](),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Last[string](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.Last[string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.Last[string](),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestLastOrElse(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.LastOrElse("D"),
		),
		"D", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.LastOrElse("D"),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.LastOrElse("D"),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.LastOrElse("D"),
		),
		"C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.LastOrElse("D"),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.LastOrElse("D"),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.LastOrElse("D"),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
