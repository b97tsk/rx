package rx_test

import (
	"testing"

	"github.com/b97tsk/rx"
	. "github.com/b97tsk/rx/internal/rxtest"
)

func TestDematerialize(t *testing.T) {
	t.Parallel()

	NewTestSuite[string](t).Case(
		rx.Pipe2(
			rx.Empty[string](),
			rx.Materialize[string](),
			rx.Dematerialize[rx.Notification[string]](),
		),
		ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Throw[string](ErrTest),
			rx.Materialize[string](),
			rx.Dematerialize[rx.Notification[string]](),
		),
		ErrTest,
	).Case(
		rx.Pipe2(
			rx.Just("A", "B", "C"),
			rx.Materialize[string](),
			rx.Dematerialize[rx.Notification[string]](),
		),
		"A", "B", "C", ErrComplete,
	).Case(
		rx.Pipe2(
			rx.Concat(
				rx.Just("A", "B", "C"),
				rx.Throw[string](ErrTest),
			),
			rx.Materialize[string](),
			rx.Dematerialize[rx.Notification[string]](),
		),
		"A", "B", "C", ErrTest,
	).Case(
		rx.Pipe1(
			rx.Empty[rx.Notification[string]](),
			rx.Dematerialize[rx.Notification[string]](),
		),
		ErrComplete,
	).Case(
		rx.Pipe1(
			rx.Throw[rx.Notification[string]](ErrTest),
			rx.Dematerialize[rx.Notification[string]](),
		),
		ErrTest,
	)
}
