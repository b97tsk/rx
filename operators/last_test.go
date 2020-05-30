package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestLast(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.Last()),
			rx.Throw(ErrTest).Pipe(operators.Last()),
			rx.Just("A").Pipe(operators.Last()),
			rx.Just("A", "B").Pipe(operators.Last()),
			rx.Concat(rx.Just("A"), rx.Throw(ErrTest)).Pipe(operators.Last()),
			rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(operators.Last()),
		},
		[][]interface{}{
			{rx.ErrEmpty},
			{ErrTest},
			{"A", rx.Completed},
			{"B", rx.Completed},
			{ErrTest},
			{ErrTest},
		},
	)
}
