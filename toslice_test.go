package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestToSlice(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		"[A B C]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Just("A"),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		"[A]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		"[]", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Oops[string](ErrTest),
			rx.ToSlice[string](),
			ToString[[]string](),
		),
		rx.ErrOops, ErrTest,
	).Case(
		rx.Pipe3(
			rx.Just("A", "B", "C"),
			rx.ToSlice[string](),
			rx.DoOnNext(func([]string) { panic(ErrTest) }),
			ToString[[]string](),
		),
		rx.ErrOops, ErrTest,
	)
}
