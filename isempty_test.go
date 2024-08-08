package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestIsEmpty(t *testing.T) {
	t.Parallel()

	NewTestSuite[bool](t).Case(
		rx.Pipe1(
			rx.Just("A", "B"),
			rx.IsEmpty[string](),
		),
		false, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Just("A"),
			rx.IsEmpty[string](),
		),
		false, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Empty[string](),
			rx.IsEmpty[string](),
		),
		true, ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[string](ErrTest),
			rx.IsEmpty[string](),
		),
		ErrTest,
	).Case(
		rx.Pipe1(
			rx.Oops[string](ErrTest),
			rx.IsEmpty[string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.IsEmpty[string](),
			rx.DoOnNext(func(bool) { panic(ErrTest) }),
		),
		rx.ErrOops, ErrTest,
	)
}
