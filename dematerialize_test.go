package rx_test

import (
	"testing"

	. "github.com/b97tsk/rx"
)

func TestOperators_Dematerialize(t *testing.T) {
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.Materialize(), operators.Dematerialize()),
			Throw(xErrTest).Pipe(operators.Materialize(), operators.Dematerialize()),
			Just("A", "B", "C").Pipe(operators.Materialize(), operators.Dematerialize()),
			Concat(Just("A", "B", "C"), Throw(xErrTest)).Pipe(operators.Materialize(), operators.Dematerialize()),
		},
		xComplete,
		xErrTest,
		"A", "B", "C", xComplete,
		"A", "B", "C", xErrTest,
	)
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.Dematerialize()),
			Throw(xErrTest).Pipe(operators.Dematerialize()),
			Just("A", "B", "C").Pipe(operators.Dematerialize()),
			Concat(Just("A", "B", "C"), Throw(xErrTest)).Pipe(operators.Dematerialize()),
		},
		xComplete,
		xErrTest,
		ErrNotNotification,
		ErrNotNotification,
	)
}
