package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFirst(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.First[string](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.First[string](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.First[string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.First[string](),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.First[string](),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.First[string](),
		),
		"A", ErrComplete,
	)
}

func TestFirstOrElse(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.FirstOrElse("D"),
		),
		"D", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.FirstOrElse("D"),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.FirstOrElse("D"),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.FirstOrElse("D"),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.FirstOrElse("D"),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Oops[string](ErrTest),
			),
			rx.FirstOrElse("D"),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.FirstOrElse("D"),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
