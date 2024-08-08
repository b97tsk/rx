package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestSingle(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe1(
			rx.Just("A"),
			rx.Single[string](),
		),
		"A", ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A", "B"),
			rx.Single[string](),
		),
		rx.ErrNotSingle,
	).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.Single[string](),
		),
		rx.ErrEmpty,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.Single[string](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.Single[string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A"),
			rx.Single[string](),
			rx.DoOnNext(func(string) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
