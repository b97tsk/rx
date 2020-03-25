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
			Throw(errTest).Pipe(operators.Materialize(), operators.Dematerialize()),
			Just("A", "B", "C").Pipe(operators.Materialize(), operators.Dematerialize()),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(operators.Materialize(), operators.Dematerialize()),
		},
		Complete,
		errTest,
		"A", "B", "C", Complete,
		"A", "B", "C", errTest,
	)
	subscribe(
		t,
		[]Observable{
			Empty().Pipe(operators.Dematerialize()),
			Throw(errTest).Pipe(operators.Dematerialize()),
			Just("A", "B", "C").Pipe(operators.Dematerialize()),
			Concat(Just("A", "B", "C"), Throw(errTest)).Pipe(operators.Dematerialize()),
		},
		Complete,
		errTest,
		ErrNotNotification,
		ErrNotNotification,
	)
}
