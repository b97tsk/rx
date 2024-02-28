package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestFromSlice(t *testing.T) {
	t.Parallel()

	s := []string{"A", "B", "C"}

	NewTestSuite[string](t).Case(
		rx.FromSlice(s),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.FromSlice(s),
			rx.Take[string](1),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.FromSlice(s),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}

func TestJust(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Just("A", "B", "C"),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.Take[string](1),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B", "C"),
			rx.OnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
