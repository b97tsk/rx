package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestDematerialize(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.Materialize(), operators.Dematerialize()),
			rx.Throw(ErrTest).Pipe(operators.Materialize(), operators.Dematerialize()),
			rx.Just("A", "B", "C").Pipe(operators.Materialize(), operators.Dematerialize()),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(operators.Materialize(), operators.Dematerialize()),
		},
		[][]interface{}{
			{rx.Completed},
			{ErrTest},
			{"A", "B", "C", rx.Completed},
			{"A", "B", "C", ErrTest},
		},
	)
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.Dematerialize()),
			rx.Throw(ErrTest).Pipe(operators.Dematerialize()),
			rx.Just("A", "B", "C").Pipe(operators.Dematerialize()),
			rx.Concat(rx.Just("A", "B", "C"), rx.Throw(ErrTest)).Pipe(operators.Dematerialize()),
		},
		[][]interface{}{
			{rx.Completed},
			{ErrTest},
			{rx.ErrNotNotification},
			{rx.ErrNotNotification},
		},
	)
}
