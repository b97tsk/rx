package operators_test

import (
	"testing"

	"github.com/b97tsk/rx"
	"github.com/b97tsk/rx/operators"
	. "github.com/b97tsk/rx/testing"
)

func TestFirst(t *testing.T) {
	SubscribeN(
		t,
		[]rx.Observable{
			rx.Empty().Pipe(operators.First()),
			rx.Throw(ErrTest).Pipe(operators.First()),
			rx.Just("A").Pipe(operators.First()),
			rx.Just("A", "B").Pipe(operators.First()),
			rx.Concat(rx.Just("A"), rx.Throw(ErrTest)).Pipe(operators.First()),
			rx.Concat(rx.Just("A", "B"), rx.Throw(ErrTest)).Pipe(operators.First()),
		},
		[][]interface{}{
			{rx.ErrEmpty},
			{ErrTest},
			{"A", rx.Completed},
			{"A", rx.Completed},
			{"A", rx.Completed},
			{"A", rx.Completed},
		},
	)
}
